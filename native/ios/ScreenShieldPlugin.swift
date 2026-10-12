import Foundation
import Capacitor
import UIKit

/**
 * ScreenShield — native screen-capture protection for the Wiyana iOS app.
 *
 * iOS reality check: Apple provides NO API to block screenshots or screen
 * recording (no FLAG_SECURE equivalent). What iOS does provide is *detection*:
 *
 * - `UIScreen.main.isCaptured` (iOS 11+) → true while the screen is being
 *   recorded OR mirrored via AirPlay. The web layer polls this during playback.
 * - `UIApplication.userDidTakeScreenshotNotification` → fires after a
 *   screenshot is taken. Forwarded as the "screenshotTaken" event so the web
 *   layer can score the violation immediately.
 *
 * Because capture can't be blocked on iOS, the response chain is:
 *   detect → pause playback instantly → warn the user → report `screen_recording`
 *   / `screenshot_attempt` to the backend (10 pts each → instant ban threshold).
 * The server-burned watermark (name + Telegram ID) remains the real deterrent:
 * any leaked recording identifies its source.
 *
 * Install: copy into ios/App/App/ and add to the Xcode target; no Podfile
 * change needed (uses only Capacitor + UIKit).
 */
@objc(ScreenShieldPlugin)
public class ScreenShieldPlugin: CAPPlugin, CAPBridgedPlugin {
    public let identifier = "ScreenShieldPlugin"
    public let jsName = "ScreenShield"
    public let pluginMethods: [CAPPluginMethod] = [
        CAPPluginMethod(name: "enableSecure", returnType: CAPPluginReturnPromise),
        CAPPluginMethod(name: "disableSecure", returnType: CAPPluginReturnPromise),
        CAPPluginMethod(name: "isCaptured", returnType: CAPPluginReturnPromise),
    ]

    private var screenshotObserver: NSObjectProtocol?
    private var captureObserver: NSObjectProtocol?

    override public func load() {
        // Screenshot taken → immediate event (the capture itself already happened;
        // the value here is instant scoring + the ban chain).
        screenshotObserver = NotificationCenter.default.addObserver(
            forName: UIApplication.userDidTakeScreenshotNotification,
            object: nil,
            queue: .main
        ) { [weak self] _ in
            self?.notifyListeners("screenshotTaken", data: ["platform": "ios"])
        }

        // Screen recording / mirroring started or stopped.
        if #available(iOS 11.0, *) {
            captureObserver = NotificationCenter.default.addObserver(
                forName: UIScreen.capturedDidChangeNotification,
                object: nil,
                queue: .main
            ) { [weak self] _ in
                let captured = UIScreen.main.isCaptured
                self?.notifyListeners("screenCaptureChanged",
                                      data: ["captured": captured, "platform": "ios"])
            }
        }
    }

    /// No-op on iOS (no FLAG_SECURE equivalent). Kept for API parity so the
    /// web layer can call it unconditionally.
    @objc func enableSecure(_ call: CAPPluginCall) {
        // Additionally hide content in the app switcher (best-effort privacy).
        call.resolve()
    }

    @objc func disableSecure(_ call: CAPPluginCall) {
        call.resolve()
    }

    /// True while the screen is being recorded or mirrored (AirPlay).
    @objc func isCaptured(_ call: CAPPluginCall) {
        var captured = false
        if #available(iOS 11.0, *) {
            captured = UIScreen.main.isCaptured
        }
        call.resolve(["captured": captured])
    }

    deinit {
        if let o = screenshotObserver { NotificationCenter.default.removeObserver(o) }
        if let o = captureObserver { NotificationCenter.default.removeObserver(o) }
    }
}
