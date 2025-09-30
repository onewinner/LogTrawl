// Wails Runtime API 类型定义
declare global {
  interface Window {
    runtime: typeof import('../wailsjs/runtime/runtime');
    backend: {
      App: typeof import('../wailsjs/go/main/App').default;
    };
  }
}

// 导出类型以供其他模块使用
export type WailsRuntime = typeof import('../wailsjs/runtime/runtime');
export type WailsBackend = {
  App: typeof import('../wailsjs/go/main/App').default;
};

export {};