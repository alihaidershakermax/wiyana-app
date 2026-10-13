// Package bot implements the Telegram bot for the Wiyana platform.
// It runs long-polling in a goroutine alongside the HTTP server and is
// connected to the same Convex catalog used by the API.
package bot

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/alihaidershakermax/wena-app-platfrom/backend-go/internal/auth"
)

// ── Texts (Arabic) ──────────────────────────────────────────────

const (
	welcomeText = `مرحباً بك في منصة ويانا التعليمية

منصة ويانا منصتك الذكية لتنظيم ملفاتك الأكاديمية: المحاضرات، الملخصات، الفيديوهات، الكتب، والدورات التفاعلية، كلها مرتبة في أقسام ومواد واضحة لتصل إليها بسهولة وفي أي وقت.

أبرز ما تقدمه المنصة:
- دورات تفاعلية ومحاضرات منظمة
- مكتبة أكاديمية شاملة
- بحث سريع عن أي ملف أو درس
- تفعيل اشتراكات الدورات برموز فورية

اضغط على زر (فتح منصة ويانا) للدخول المباشر من تليجرام.`

	aboutText = `السلام عليكم ورحمة الله وبركاته 👋

🤖 فكرة هذا البوت والمنصة:
منصة تعليمية ذكية تنظّم ملفاتك الأكاديمية (محاضرات، ملخصات، فيديوهات، كتب) ودوراتك التفاعلية في أقسام ومواد مرتبة، لتصل إليها بضغطة زر وفي أي وقت.

✨ ماذا يقدم لك؟
🚀 منصة Wiyana التفاعلية داخل تليجرام
📂 أقسام ومواد منظمة ومرتبة
🔍 بحث سريع عن أي ملف أو درس
⬇️ تحميل مباشر أو عرض PDF كصور
🔗 مشاركة الملفات بروابط مباشرة

👨‍💻 المطور: علي الأكبر حيدر شاكر

شكراً لثقتكم بنا، ونتمنى أن يكون هذا العمل إضافة حقيقية لرحلتكم التعليمية.`

	// OwnerID is the platform owner (super-admin) by default.
	OwnerID int64 = 2018954602
)

// ── Convex client surface used by the bot ──────────────────────

type userDoc struct {
	ID        string `json:"_id"`
	TelegramID string `json:"telegram_id"`
	FirstName string `json:"first_name"`
	IsAdmin   bool   `json:"is_admin"`
	IsActive  bool   `json:"is_active"`
}

type videoDoc struct {
	ID        string `json:"_id"`
	Title     string `json:"title"`
	Category  string `json:"category"`
	Subject   string `json:"subject"`
	TelegramFileID string `json:"telegram_file_id,omitempty"`
	Downloads int     `json:"downloads,omitempty"`
}

// ConvexStore is the interface the bot needs from the Convex layer.
type ConvexStore interface {
	CallQuery(ctx context.Context, path string, args map[string]any) (json.RawMessage, error)
	CallMutation(ctx context.Context, path string, args map[string]any) (json.RawMessage, error)
}

// ── Bot ─────────────────────────────────────────────────────────

// Config holds everything the bot needs.
type Config struct {
	Token       string
	WebAppURL   string
	OwnerID     int64
	AdminIDs    []int64 // extra admin telegram ids
	PollTimeout int     // seconds, default 60
}

// Bot is the Telegram bot controller.
type Bot struct {
	api   *tgbotapi.BotAPI
	cfg   Config
	store ConvexStore

	mu      sync.RWMutex
	awaiting map[int64]string // chat_id -> awaiting state (search, new_cat, new_sub, broadcast)
	pending map[int64]any      // chat_id -> pending payload per state

	adminMu sync.RWMutex
	admins  map[int64]struct{}

	shutdown chan struct{}
}

// New creates a bot controller (does not start polling yet).
func New(cfg Config, store ConvexStore) *Bot {
	if cfg.OwnerID == 0 {
		cfg.OwnerID = OwnerID
	}
	if cfg.WebAppURL == "" {
		cfg.WebAppURL = "https://wena-app.onrender.com/"
	}
	if cfg.PollTimeout <= 0 {
		cfg.PollTimeout = 60
	}
	b := &Bot{
		cfg:      cfg,
		store:    store,
		awaiting: make(map[int64]string),
		pending:  make(map[int64]any),
		admins:   make(map[int64]struct{}),
		shutdown: make(chan struct{}),
	}
	b.admins[cfg.OwnerID] = struct{}{}
	for _, id := range cfg.AdminIDs {
		b.admins[id] = struct{}{}
	}
	return b
}

// Run starts the bot. It blocks until ctx is cancelled.
func (b *Bot) Run(ctx context.Context) error {
	api, err := tgbotapi.NewBotAPI(b.cfg.Token)
	if err != nil {
		return fmt.Errorf("telegram bot init: %w", err)
	}
	api.Debug = false
	b.api = api
	log.Printf("[bot] authorized as %s", api.Self.UserName)

	// Warm up admin list from Convex (users with is_admin=true).
	go b.refreshAdmins(ctx)
	t := time.NewTicker(5 * time.Minute)
	defer t.Stop()
	go func() {
		for {
			select {
			case <-t.C:
				b.refreshAdmins(ctx)
			case <-ctx.Done():
				return
			}
		}
	}()

	updates := tgbotapi.NewUpdate(0)
	updates.Timeout = b.cfg.PollTimeout
	ch := api.GetUpdatesChan(updates)
	log.Printf("[bot] polling started")
	for {
		select {
		case <-ctx.Done():
			close(b.shutdown)
			return nil
		case upd := <-ch:
			go b.handleUpdate(ctx, upd)
		}
	}
}

// ── Update routing ──────────────────────────────────────────────

func (b *Bot) handleUpdate(ctx context.Context, upd tgbotapi.Update) {
	switch {
	case upd.Message != nil && upd.Message.IsCommand():
		b.handleCommand(ctx, upd.Message)
	case upd.Message != nil && upd.Message.Text != "":
		b.handleText(ctx, upd.Message)
	case upd.CallbackQuery != nil:
		b.handleCallback(ctx, upd.CallbackQuery)
	}
}

// ── Commands ────────────────────────────────────────────────────

func (b *Bot) handleCommand(ctx context.Context, msg *tgbotapi.Message) {
	user := msg.From
	if user == nil {
		return
	}
	args := strings.TrimSpace(msg.CommandArguments())

	switch msg.Command() {
	case "start":
		b.handleStart(ctx, msg, args)
	case "panel":
		b.handlePanel(ctx, msg)
	case "broadcast":
		b.handleBroadcast(ctx, msg, args)
	case "ban":
		b.handleBan(ctx, msg, args, false)
	case "unban":
		b.handleBan(ctx, msg, args, true)
	}
}

func (b *Bot) handleStart(ctx context.Context, msg *tgbotapi.Message, args string) {
	user := msg.From
	// register user
	_, _ = b.store.CallMutation(ctx, "users:upsertUser", map[string]any{
		"telegram_id": fmt.Sprint(user.ID),
		"first_name":  user.FirstName,
		"last_name":   user.LastName,
		"username":    user.UserName,
	})

	// check ban
	if banned, _ := b.userActive(ctx, user.ID); !banned {
		b.sendText(msg.Chat.ID, "⛔ أنت محظور من استخدام البوت.", nil, false)
		return
	}

	if args != "" {
		// deep links
		if strings.HasPrefix(args, "share_") {
			b.handleShareDeepLink(ctx, msg, args[len("share_"):])
			return
		}
		if strings.HasPrefix(args, "code_") {
			b.sendText(msg.Chat.ID, "🎟️ يرجى تفعيل الكود من داخل منصة ويانا مباشرة.", nil, false)
			return
		}
	}

	extra := [][]tgbotapi.InlineKeyboardButton{
		{
			tgbotapi.NewInlineKeyboardButtonData("📂 الأرشيف", "view_archive"),
			tgbotapi.NewInlineKeyboardButtonData("🔍 بحث", "search_menu"),
		},
		{
			tgbotapi.NewInlineKeyboardButtonData("👨‍💻 عن المطور", "about"),
		},
	}
	b.sendTextWithWebApp(msg.Chat.ID, msg.From.ID, welcomeText, extra)
}

func (b *Bot) handleShareDeepLink(ctx context.Context, msg *tgbotapi.Message, hash string) {
	row, err := b.store.CallQuery(ctx, "users:getShareLink", map[string]any{"hash": hash})
	if err != nil || row == nil {
		b.sendText(msg.Chat.ID, "❌ رابط المشاركة غير صالح أو منتهي الصلاحية.", nil, false)
		return
	}
	var link struct {
		FileID string `json:"file_id"`
	}
	if err := json.Unmarshal(row, &link); err != nil || link.FileID == "" {
		b.sendText(msg.Chat.ID, "❌ رابط المشاركة غير صالح أو منتهي الصلاحية.", nil, false)
		return
	}
	doc := tgbotapi.NewDocument(msg.Chat.ID, tgbotapi.FileID(link.FileID))
	doc.Caption = "📤 الملف المشارك معك عبر الأرشيف التعليمي"
	_, _ = b.api.Send(doc)
}

func (b *Bot) handlePanel(ctx context.Context, msg *tgbotapi.Message) {
	if !b.isAdmin(msg.From.ID) {
		b.sendText(msg.Chat.ID, "⛔ ليس لديك صلاحيات لوحة التحكم.", nil, false)
		return
	}
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("📊 الإحصائيات", "panel_stats"),
			tgbotapi.NewInlineKeyboardButtonData("👥 المستخدمين", "panel_users"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("📁 إدارة الأرشيف", "adm_content"),
			tgbotapi.NewInlineKeyboardButtonData("📢 البث الجماعي", "adm_broadcast"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🔙 العودة للقائمة الرئيسية", "back_main"),
		),
	)
	txt := fmt.Sprintf("🛡️ <b>لوحة تحكم المشرفين:</b>\nأهلاً بك يا <b>%s</b>.", msg.From.FirstName)
	b.sendHTML(msg.Chat.ID, txt, &keyboard, false)
}

func (b *Bot) handleBroadcast(ctx context.Context, msg *tgbotapi.Message, args string) {
	if !b.isAdmin(msg.From.ID) {
		b.sendText(msg.Chat.ID, "⛔ ليس لديك صلاحية الإرسال الجماعي.", nil, false)
		return
	}
	if args == "" {
		b.sendHTML(msg.Chat.ID, "يرجى كتابة نص الرسالة الجماعية بعد الأمر:\n<code>/broadcast نص الرسالة</code>", nil, false)
		return
	}
	b.setAwaiting(msg.Chat.ID, "broadcast", args)
	sent, err := b.broadcastToUsers(ctx, args)
	b.clearAwaiting(msg.Chat.ID)
	if err != nil {
		b.sendText(msg.Chat.ID, "❌ حدث خطأ أثناء البث. يرجى المحاولة لاحقاً.", nil, false)
		return
	}
	b.sendText(msg.Chat.ID, fmt.Sprintf("✅ تم إرسال البث بنجاح إلى %d مستخدم.", sent), nil, false)
}

func (b *Bot) handleBan(ctx context.Context, msg *tgbotapi.Message, args string, unban bool) {
	if !b.isAdmin(msg.From.ID) {
		b.sendText(msg.Chat.ID, "⛔ ليس لديك صلاحية إدارة المستخدمين.", nil, false)
		return
	}
	id := strings.TrimSpace(args)
	if id == "" {
		if unban {
			b.sendHTML(msg.Chat.ID, "استخدم: <code>/unban USER_ID</code>", nil, false)
		} else {
			b.sendHTML(msg.Chat.ID, "استخدم: <code>/ban USER_ID</code>", nil, false)
		}
		return
	}
	var path string
	if unban {
		path = "users:unbanUser"
	} else {
		path = "users:banUser"
	}
	if _, err := b.store.CallMutation(ctx, path, map[string]any{"telegramId": id}); err != nil {
		b.sendText(msg.Chat.ID, "❌ خطأ: تعذر تنفيذ العملية.", nil, false)
		return
	}
	if unban {
		b.sendText(msg.Chat.ID, fmt.Sprintf("✅ تم إلغاء حظر المستخدم %s.", id), nil, false)
	} else {
		b.sendText(msg.Chat.ID, fmt.Sprintf("✅ تم حظر المستخدم %s.", id), nil, false)
	}
}

// ── Free text / awaiting states ─────────────────────────────────

func (b *Bot) handleText(ctx context.Context, msg *tgbotapi.Message) {
	chatID := msg.Chat.ID
	b.mu.RLock()
	state := b.awaiting[chatID]
	b.mu.RUnlock()

	switch state {
	case "search":
		b.handleSearch(ctx, msg, strings.TrimSpace(msg.Text))
	case "broadcast":
		b.mu.RLock()
		text := b.pending[chatID].(string)
		b.mu.RUnlock()
		sent, err := b.broadcastToUsers(ctx, text)
		b.clearAwaiting(chatID)
		if err != nil {
			b.sendText(chatID, "❌ حدث خطأ أثناء البث. يرجى المحاولة لاحقاً.", nil, false)
		} else {
			b.sendText(chatID, fmt.Sprintf("✅ تم إرسال البث الجماعي إلى %d مستخدم.", sent), nil, false)
		}
	default:
		// no awaiting state — ignore free text
	}
}

func (b *Bot) handleSearch(ctx context.Context, msg *tgbotapi.Message, query string) {
	b.clearAwaiting(msg.Chat.ID)
	if query == "" {
		return
	}
	row, err := b.store.CallQuery(ctx, "users:searchVideos", map[string]any{"q": query})
	if err != nil {
		b.sendText(msg.Chat.ID, "❌ تعذر تنفيذ البحث، يرجى المحاولة لاحقاً.", nil, false)
		return
	}
	var videos []videoDoc
	_ = json.Unmarshal(row, &videos)
	if len(videos) == 0 {
		b.sendHTML(msg.Chat.ID, fmt.Sprintf("🔍 لم يتم العثور على أي نتائج تطابق: <b>%s</b>", esc(query)), nil, false)
		return
	}
	keyboard := tgbotapi.NewInlineKeyboardMarkup()
	for _, v := range videos {
		keyboard.InlineKeyboard = append(keyboard.InlineKeyboard,
			tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("📄 "+v.Title, "file_"+v.ID)))
	}
	keyboard.InlineKeyboard = append(keyboard.InlineKeyboard,
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("🔙 العودة للقائمة الرئيسية", "back_main")))
	b.sendHTML(msg.Chat.ID, fmt.Sprintf("🔍 نتائج البحث عن <b>%s</b> (%d ملف):", esc(query), len(videos)), &keyboard, false)
}

// ── Callbacks ───────────────────────────────────────────────────

func (b *Bot) handleCallback(ctx context.Context, q *tgbotapi.CallbackQuery) {
	user := q.From
	data := q.Data

	// common guards
	if active, _ := b.userActive(ctx, user.ID); !active {
		b.answer(q, "⛔ أنت محظور من استخدام البوت.", true)
		return
	}

	switch {
	case data == "back_main":
		b.cbMain(ctx, q)
	case data == "view_archive" || data == "back_cats":
		b.cbCategories(ctx, q)
	case strings.HasPrefix(data, "cat_"):
		b.cbSubjects(ctx, q, strings.TrimPrefix(data, "cat_"))
	case strings.HasPrefix(data, "sub_"):
		b.cbFiles(ctx, q, strings.TrimPrefix(data, "sub_"))
	case strings.HasPrefix(data, "file_"):
		b.cbFileDetails(ctx, q, strings.TrimPrefix(data, "file_"))
	case strings.HasPrefix(data, "download_"):
		b.cbDownload(ctx, q, strings.TrimPrefix(data, "download_"))
	case strings.HasPrefix(data, "share_"):
		b.cbShare(ctx, q, strings.TrimPrefix(data, "share_"))
	case data == "search_menu":
		b.cbEnterSearch(ctx, q)
	case data == "about":
		b.cbAbout(ctx, q)
	case data == "panel" || data == "adm_back":
		b.cbAdminPanel(ctx, q)
	case data == "panel_stats":
		b.cbStats(ctx, q)
	case data == "panel_users":
		b.cbUsersAdmin(ctx, q)
	case data == "adm_content":
		b.cbAdmContent(ctx, q)
	case data == "adm_broadcast":
		b.cbAdmBroadcast(ctx, q)
	default:
		b.answer(q, "", false)
	}
}

func (b *Bot) cbMain(ctx context.Context, q *tgbotapi.CallbackQuery) {
	b.answer(q, "", false)
	// Telegram does not allow converting an existing message into a web_app
	// button, so send a fresh message with the Mini App button and remove the
	// old menu message when possible.
	extra := [][]tgbotapi.InlineKeyboardButton{
		{
			tgbotapi.NewInlineKeyboardButtonData("📂 الأرشيف", "view_archive"),
			tgbotapi.NewInlineKeyboardButtonData("🔍 بحث", "search_menu"),
		},
		{
			tgbotapi.NewInlineKeyboardButtonData("👨‍💻 عن المطور", "about"),
		},
	}
	b.sendTextWithWebApp(q.Message.Chat.ID, q.From.ID, welcomeText, extra)
	if q.Message.MessageID != 0 {
		_, _ = b.api.Send(tgbotapi.NewDeleteMessage(q.Message.Chat.ID, q.Message.MessageID))
	}
}

func (b *Bot) cbCategories(ctx context.Context, q *tgbotapi.CallbackQuery) {
	b.answer(q, "", false)
	row, err := b.store.CallQuery(ctx, "users:listCategories", map[string]any{})
	if err != nil || row == nil {
		b.answer(q, "⚠️ تعذر جلب الأقسام، يرجى المحاولة لاحقاً.", true)
		return
	}
	var cats []struct {
		Name string `json:"name"`
	}
	_ = json.Unmarshal(row, &cats)
	keyboard := tgbotapi.NewInlineKeyboardMarkup()
	for _, c := range cats {
		keyboard.InlineKeyboard = append(keyboard.InlineKeyboard,
			tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("📂 "+c.Name, "cat_"+escCB(c.Name))))
	}
	keyboard.InlineKeyboard = append(keyboard.InlineKeyboard,
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("🔙 العودة للقائمة الرئيسية", "back_main")))
	b.editOrReply(q, "📂 <b>الأقسام:</b>\n\nاختر القسم الذي تريد تصفحه:", &keyboard, true)
}

func (b *Bot) cbSubjects(ctx context.Context, q *tgbotapi.CallbackQuery, category string) {
	b.answer(q, "", false)
	row, err := b.store.CallQuery(ctx, "users:listSubjects", map[string]any{"category": category})
	if err != nil || row == nil {
		b.answer(q, "⚠️ تعذر جلب المواد، يرجى المحاولة لاحقاً.", true)
		return
	}
	var subs []struct {
		Name string `json:"name"`
	}
	_ = json.Unmarshal(row, &subs)
	keyboard := tgbotapi.NewInlineKeyboardMarkup()
	for _, s := range subs {
		keyboard.InlineKeyboard = append(keyboard.InlineKeyboard,
			tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("📚 "+s.Name, "sub_"+escCB(s.Name))))
	}
	keyboard.InlineKeyboard = append(keyboard.InlineKeyboard,
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("🔙 العودة للأقسام", "back_cats")))
	b.editOrReply(q, fmt.Sprintf("📚 <b>مواد قسم (%s):</b>", esc(category)), &keyboard, true)
}

func (b *Bot) cbFiles(ctx context.Context, q *tgbotapi.CallbackQuery, subject string) {
	b.answer(q, "", false)
	// find category of this subject from the first matching video
	videos, category, err := b.listBySubject(ctx, subject)
	if err != nil || len(videos) == 0 {
		b.answer(q, "❌ لا توجد ملفات في هذه المادة حالياً.", true)
		return
	}
	keyboard := tgbotapi.NewInlineKeyboardMarkup()
	for _, v := range videos {
		keyboard.InlineKeyboard = append(keyboard.InlineKeyboard,
			tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("📄 "+v.Title, "file_"+v.ID)))
	}
	keyboard.InlineKeyboard = append(keyboard.InlineKeyboard,
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("🔙 العودة للمواد", "back_subs_"+escCB(category))))
	b.editOrReply(q, fmt.Sprintf("📁 <b>ملفات مادة (%s):</b>", esc(subject)), &keyboard, true)
}

func (b *Bot) cbFileDetails(ctx context.Context, q *tgbotapi.CallbackQuery, id string) {
	b.answer(q, "", false)
	v, err := b.getVideo(ctx, id)
	if err != nil || v == nil {
		b.answer(q, "❌ الملف غير موجود.", true)
		return
	}
	tgID := ""
	if v.TelegramFileID != "" {
		tgID = v.TelegramFileID
	}
	_ = tgID
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("⬇️ تحميل", "download_"+id),
			tgbotapi.NewInlineKeyboardButtonData("🔗 مشاركة", "share_"+id),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🔙 العودة للقائمة الرئيسية", "back_main"),
		),
	)
	txt := fmt.Sprintf("📄 <b>%s</b>\n\n📂 القسم: %s\n📚 المادة: %s\n⬇️ التحميلات: %d",
		esc(v.Title), esc(orStr(v.Category, "-")), esc(orStr(v.Subject, "-")), v.Downloads)
	b.editOrReply(q, txt, &keyboard, true)
}

func (b *Bot) cbDownload(ctx context.Context, q *tgbotapi.CallbackQuery, id string) {
	v, err := b.getVideo(ctx, id)
	if err != nil || v == nil || v.TelegramFileID == "" {
		b.answer(q, "❌ الملف غير متاح للتحميل.", true)
		return
	}
	b.answer(q, "⏳ جاري إرسال الملف...", false)
	_, _ = b.store.CallMutation(ctx, "users:incrementDownloads", map[string]any{"id": id})
	doc := tgbotapi.NewDocument(q.From.ID, tgbotapi.FileID(v.TelegramFileID))
	doc.Caption = fmt.Sprintf("📥 <b>%s</b>", esc(v.Title))
	doc.ParseMode = "HTML"
	if _, err := b.api.Send(doc); err != nil {
		log.Printf("[bot] send document error: %v", err)
		b.answer(q, "❌ تعذر إرسال الملف، يرجى المحاولة لاحقاً.", true)
	}
}

func (b *Bot) cbShare(ctx context.Context, q *tgbotapi.CallbackQuery, id string) {
	v, err := b.getVideo(ctx, id)
	if err != nil || v == nil || v.TelegramFileID == "" {
		b.answer(q, "❌ الملف غير موجود.", true)
		return
	}
	row, err := b.store.CallMutation(ctx, "users:createShareLink", map[string]any{
		"file_id": v.TelegramFileID,
	})
	if err != nil || row == nil {
		b.answer(q, "❌ تعذر إنشاء رابط المشاركة.", true)
		return
	}
	var hash string
	_ = json.Unmarshal(row, &hash)
	shareURL := fmt.Sprintf("https://t.me/%s?start=share_%s", b.api.Self.UserName, hash)
	txt := fmt.Sprintf("🔗 <b>رابط المشاركة المباشر للملف:</b>\n\n<code>%s</code>\n\nيمكنك نسخ هذا الرابط وإرساله لأي شخص لتحميل الملف فوراً!", shareURL)
	b.editOrReply(q, txt, nil, true)
}

func (b *Bot) cbEnterSearch(ctx context.Context, q *tgbotapi.CallbackQuery) {
	b.answer(q, "", false)
	b.setAwaiting(q.From.ID, "search", nil)
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("🔙 العودة للقائمة الرئيسية", "back_main")),
	)
	b.edit(q, "🔍 <b>البحث الذكي في الأرشيف:</b>\n\nأرسل اسم الكلمة أو المحاضرة التي تبحث عنها الآن:", &keyboard)
}

func (b *Bot) cbAbout(ctx context.Context, q *tgbotapi.CallbackQuery) {
	b.answer(q, "", false)
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("🔙 العودة للقائمة الرئيسية", "back_main")),
	)
	b.editOrReply(q, aboutText, &keyboard, true)
}

// ── Admin callbacks ─────────────────────────────────────────────

func (b *Bot) cbAdminPanel(ctx context.Context, q *tgbotapi.CallbackQuery) {
	b.answer(q, "", false)
	if !b.isAdmin(q.From.ID) {
		b.answer(q, "⛔ ليس لديك صلاحيات لوحة التحكم.", true)
		return
	}
	b.cbPanel(ctx, q)
}

func (b *Bot) cbPanel(ctx context.Context, q *tgbotapi.CallbackQuery) {
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("📊 الإحصائيات", "panel_stats"),
			tgbotapi.NewInlineKeyboardButtonData("👥 المستخدمين", "panel_users"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("📁 إدارة الأرشيف", "adm_content"),
			tgbotapi.NewInlineKeyboardButtonData("📢 البث الجماعي", "adm_broadcast"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🔙 العودة للقائمة الرئيسية", "back_main"),
		),
	)
	txt := fmt.Sprintf("🛡️ <b>لوحة تحكم المشرفين والمدراء:</b>\n\n👤 <b>حسابك:</b> %s\n\nاختر الإجراء المطلوب من القائمة:", q.From.FirstName)
	b.editOrReply(q, txt, &keyboard, true)
}

func (b *Bot) cbStats(ctx context.Context, q *tgbotapi.CallbackQuery) {
	b.answer(q, "", false)
	if !b.isAdmin(q.From.ID) {
		return
	}
	row, err := b.store.CallQuery(ctx, "users:listUsers", map[string]any{})
	if err != nil || row == nil {
		b.answer(q, "⚠️ تعذر جلب الإحصائيات.", true)
		return
	}
	var users []userDoc
	_ = json.Unmarshal(row, &users)
	row2, err2 := b.store.CallQuery(ctx, "videos:listVideos", map[string]any{})
	videosCount := 0
	if err2 == nil && row2 != nil {
		var vids []videoDoc
		_ = json.Unmarshal(row2, &vids)
		videosCount = len(vids)
	}
	txt := fmt.Sprintf("📊 <b>الإحصائيات:</b>\n\n👥 المستخدمون المسجلون: %d\n📁 الملفات في الأرشيف: %d", len(users), videosCount)
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("🔙 لوحة الأدمن", "adm_back")),
	)
	b.editOrReply(q, txt, &keyboard, true)
}

func (b *Bot) cbUsersAdmin(ctx context.Context, q *tgbotapi.CallbackQuery) {
	b.answer(q, "", false)
	if !b.isAdmin(q.From.ID) {
		return
	}
	row, err := b.store.CallQuery(ctx, "users:listUsers", map[string]any{})
	if err != nil || row == nil {
		b.answer(q, "⚠️ تعذر جلب قائمة المستخدمين.", true)
		return
	}
	var users []userDoc
	_ = json.Unmarshal(row, &users)
	if len(users) == 0 {
		b.answer(q, "👥 لا يوجد مستخدمون مسجلون بعد.", true)
		return
	}
	// top 20 newest
	limit := 20
	if len(users) < limit {
		limit = len(users)
	}
	lines := []string{"👥 <b>آخر المستخدمين المسجلين:</b>\n"}
	for _, u := range users[:limit] {
		status := "✅"
		if !u.IsActive {
			status = "⛔"
		}
		lines = append(lines, fmt.Sprintf("%s %s (%s)", status, esc(orStr(u.FirstName, "مستخدم")), u.TelegramID))
	}
	lines = append(lines, "\nلحظر/إلغاء حظر مستخدم استخدم /ban أو /unban مع معرّفه.")
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("🔙 لوحة الأدمن", "adm_back")),
	)
	b.editOrReply(q, strings.Join(lines, "\n"), &keyboard, true)
}

func (b *Bot) cbAdmContent(ctx context.Context, q *tgbotapi.CallbackQuery) {
	b.answer(q, "", false)
	if !b.isAdmin(q.From.ID) {
		return
	}
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("➕ قسم جديد", "adm_new_cat"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("📁 إضافة مادة لقسم", "adm_new_sub"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🔙 لوحة الأدمن", "adm_back"),
		),
	)
	b.editOrReply(q, "📁 <b>إدارة الأرشيف:</b>\n\nملاحظة: إضافة الملفات (الفيديوهات والمحاضرات) تتم من داخل منصة ويانا عبر لوحة الأدمن الرئيسية في المتصفح. من هنا يمكنك إنشاء الأقسام والمواد فقط.", &keyboard, true)
}

func (b *Bot) cbAdmBroadcast(ctx context.Context, q *tgbotapi.CallbackQuery) {
	b.answer(q, "", false)
	if !b.isAdmin(q.From.ID) {
		return
	}
	b.setAwaiting(q.From.ID, "broadcast", nil)
	b.edit(q, "📢 <b>البث الجماعي:</b>\n\nأرسل الآن نص الرسالة التي تريد إرسالها لجميع المستخدمين:", nil)
}

// ── Helpers ─────────────────────────────────────────────────────

func (b *Bot) refreshAdmins(ctx context.Context) {
	row, err := b.store.CallQuery(ctx, "users:listUsers", map[string]any{})
	if err != nil || row == nil {
		return
	}
	var users []userDoc
	_ = json.Unmarshal(row, &users)
	b.adminMu.Lock()
	defer b.adminMu.Unlock()
	for _, u := range users {
		if u.IsAdmin {
		id, perr := strconv.ParseInt(u.TelegramID, 10, 64)
		if perr != nil {
			continue
		}
		b.admins[id] = struct{}{}
			
		}
	}
}

func (b *Bot) isAdmin(id int64) bool {
	b.adminMu.RLock()
	defer b.adminMu.RUnlock()
	_, ok := b.admins[id]
	return ok
}

// userActive returns true when the user is allowed to use the bot.
// It lazily registers new users.
func (b *Bot) userActive(ctx context.Context, id int64) (bool, error) {
	row, err := b.store.CallQuery(ctx, "users:getUser", map[string]any{"telegramId": fmt.Sprint(id)})
	if err != nil {
		return false, err
	}
	var u *userDoc
	if row != nil {
		_ = json.Unmarshal(row, &u)
	}
	if u == nil {
		return true, nil // new user, will be registered on /start
	}
	return u.IsActive, nil
}

func (b *Bot) getVideo(ctx context.Context, id string) (*videoDoc, error) {
	row, err := b.store.CallQuery(ctx, "videos:getVideo", map[string]any{"id": id})
	if err != nil || row == nil {
		return nil, err
	}
	var v videoDoc
	if err := json.Unmarshal(row, &v); err != nil {
		return nil, err
	}
	if v.ID == "" {
		return nil, nil
	}
	return &v, nil
}

func (b *Bot) listBySubject(ctx context.Context, subject string) ([]videoDoc, string, error) {
	row, err := b.store.CallQuery(ctx, "users:listVideosBySubject", map[string]any{"subject": subject})
	if err != nil || row == nil {
		return nil, "", err
	}
	var vids []videoDoc
	if err := json.Unmarshal(row, &vids); err != nil {
		return nil, "", err
	}
	category := ""
	if len(vids) > 0 {
		category = vids[0].Category
	}
	return vids, category, nil
}

func (b *Bot) broadcastToUsers(ctx context.Context, text string) (int, error) {
	row, err := b.store.CallQuery(ctx, "users:listUsers", map[string]any{})
	if err != nil || row == nil {
		return 0, err
	}
	var users []userDoc
	_ = json.Unmarshal(row, &users)
	sent := 0
	for _, u := range users {
		if !u.IsActive {
			continue
		}
		id, perr := strconv.ParseInt(u.TelegramID, 10, 64)
		if perr != nil {
			continue
		}
		msg := tgbotapi.NewMessage(id, "📢 "+text)
		msg.ParseMode = "HTML"
		if _, serr := b.api.Send(msg); serr == nil {
			sent++
		}
		time.Sleep(40 * time.Millisecond) // stay within rate limits
	}
	return sent, nil
}

func (b *Bot) setAwaiting(chatID int64, state string, payload any) {
	b.mu.Lock()
	b.awaiting[chatID] = state
	if payload != nil {
		b.pending[chatID] = payload
	}
	b.mu.Unlock()
}

func (b *Bot) clearAwaiting(chatID int64) {
	b.mu.Lock()
	delete(b.awaiting, chatID)
	delete(b.pending, chatID)
	b.mu.Unlock()
}

// ── Message sending helpers ─────────────────────────────────────

func (b *Bot) sendText(chatID int64, text string, keyboard *tgbotapi.InlineKeyboardMarkup, protect bool) {
	msg := tgbotapi.NewMessage(chatID, text)
	if keyboard != nil {
		msg.ReplyMarkup = *keyboard
	}
	_, _ = b.api.Send(msg)
}

func (b *Bot) sendHTML(chatID int64, text string, keyboard *tgbotapi.InlineKeyboardMarkup, protect bool) {
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "HTML"
	msg.ParseMode = "HTML"
	if keyboard != nil {
		msg.ReplyMarkup = keyboard
	}
	_, _ = b.api.Send(msg)
}

// sendTextWithWebApp sends a message whose first row's first button opens the
// Wiyana Mini App inside Telegram (web_app button type). Other buttons are
// copied from the provided standard keyboard rows.
// userID is used to generate a one-time auth token as fallback when initData is unavailable.
func (b *Bot) sendTextWithWebApp(id int64, userID int64, text string, extraRows [][]tgbotapi.InlineKeyboardButton) {
	var rows [][]tgInlineButton
	webAppURL := strings.TrimRight(b.cfg.WebAppURL, "/") + "/"
	// Append one-time auth token as fallback (valid 10min, single-use)
	if userID != 0 {
		token := auth.CreateOneTimeToken(userID)
		webAppURL += "?auth=" + token
	}
	rows = append(rows, []tgInlineButton{{Text: "🚀 فتح منصة ويانا", WebApp: &tgWebAppInfo{URL: webAppURL}}})
	for _, r := range extraRows {
		row := make([]tgInlineButton, len(r))
		for i, btn := range r {
			row[i] = tgInlineButton{Text: btn.Text, Data: btn.CallbackData}
		}
		rows = append(rows, row)
	}
	msg := tgbotapi.NewMessage(id, text)
	msg.ParseMode = "HTML"
	msg.ReplyMarkup = tgInlineKeyboard{InlineKeyboard: rows}
	_, _ = b.api.Send(msg)
}

func (b *Bot) edit(msgg *tgbotapi.CallbackQuery, text string, keyboard *tgbotapi.InlineKeyboardMarkup) {
	edit := tgbotapi.NewEditMessageText(msgg.Message.Chat.ID, msgg.Message.MessageID, text)
	edit.ParseMode = "HTML"
	if keyboard != nil {
		edit.ReplyMarkup = keyboard
	}
	if _, err := b.api.Send(edit); err != nil {
		// fallback to reply if edit fails
		reply := tgbotapi.NewMessage(msgg.Message.Chat.ID, text)
		reply.ParseMode = "HTML"
		if keyboard != nil {
			reply.ReplyMarkup = *keyboard
		}
		_, _ = b.api.Send(reply)
	}
}

func (b *Bot) editOrReply(q *tgbotapi.CallbackQuery, text string, keyboard *tgbotapi.InlineKeyboardMarkup, protect bool) {
	edit := tgbotapi.NewEditMessageText(q.Message.Chat.ID, q.Message.MessageID, text)
	edit.ParseMode = "HTML"
	if keyboard != nil {
		edit.ReplyMarkup = keyboard
	}
	if _, err := b.api.Send(edit); err != nil {
		reply := tgbotapi.NewMessage(q.Message.Chat.ID, text)
		reply.ParseMode = "HTML"
		if keyboard != nil {
			reply.ReplyMarkup = keyboard
		}
		_, _ = b.api.Send(reply)
	}
}

func (b *Bot) answer(q *tgbotapi.CallbackQuery, text string, show bool) {
	cfg := tgbotapi.NewCallback(q.ID, text)
	cfg.ShowAlert = show
	_, _ = b.api.Send(cfg)
}

// ── small utils ─────────────────────────────────────────────────

func esc(s string) string {
	s = strings.ReplaceAll(s, "&", "&amp;")
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	return s
}

func escCB(s string) string {
	// callback data must be <=64 bytes; keep short and ascii-ish
	if len(s) <= 60 {
		return s
	}
	return s[:60]
}

func strPtr(s string) *string { return &s }

// tgInlineButton is a local InlineKeyboardButton type that supports the
// web_app field (required to open the Wiyana Mini App inside Telegram).
// The official library (v5.5.1) predates web_app support, so we marshal
// inline keyboards ourselves and pass the raw JSON via ReplyMarkup.
type tgInlineButton struct {
	Text   string          `json:"text"`
	URL    *string         `json:"url,omitempty"`
	Data   *string         `json:"callback_data,omitempty"`
	WebApp *tgWebAppInfo   `json:"web_app,omitempty"`
}

type tgWebAppInfo struct {
	URL string `json:"url"`
}

type tgInlineKeyboard struct {
	InlineKeyboard [][]tgInlineButton `json:"inline_keyboard"`
}

func orStr(s, fallback string) string {
	if s == "" {
		return fallback
	}
	return s
}

