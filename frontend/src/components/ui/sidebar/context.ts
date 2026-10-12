import { ref, inject, provide, computed, readonly, type InjectionKey, type Ref, type ComputedRef } from "vue";

export interface SidebarContext {
  open: Ref<boolean>;
  setOpen: (value: boolean) => void;
  openMobile: Ref<boolean>;
  setOpenMobile: (value: boolean) => void;
  isMobile: Ref<boolean>;
  toggleSidebar: () => void;
  state: ComputedRef<"expanded" | "collapsed">;
}

export const SIDEBAR_INJECTION_KEY: InjectionKey<SidebarContext> = Symbol("sidebar_context");

export function useSidebar() {
  const context = inject(SIDEBAR_INJECTION_KEY);
  if (!context) {
    throw new Error("useSidebar must be used within a SidebarProvider");
  }
  return context;
}

