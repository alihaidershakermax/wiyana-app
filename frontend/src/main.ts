import { createApp } from "vue";
import { createPinia } from "pinia";
import { VueQueryPlugin } from "@tanstack/vue-query";
import { MotionPlugin } from "@vueuse/motion";
import App from "./App.vue";
import router from "./router";
import "./styles.css";

const app = createApp(App);

// Global error handler to catch runtime Vue component errors
app.config.errorHandler = (err, instance, info) => {
  console.error("[Vue Error Handler]", err);
  console.error("[Component]", instance);
  console.error("[Info]", info);
};

app.use(createPinia());
app.use(router);
app.use(VueQueryPlugin);
app.use(MotionPlugin);

app.mount("#app");
