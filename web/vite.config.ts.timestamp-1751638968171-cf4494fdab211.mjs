// vite.config.ts
import { reactRouter } from "file:///home/asim/go/src/github.com/asim/reminder/web/node_modules/.pnpm/@react-router+dev@7.5.1_@react-router+serve@7.5.1_react-router@7.5.1_react-dom@19.1.0_r_461469edd158ee3394e305a5ce94f09f/node_modules/@react-router/dev/dist/vite.js";
import tailwindcss from "file:///home/asim/go/src/github.com/asim/reminder/web/node_modules/.pnpm/@tailwindcss+vite@4.1.4_vite@5.4.18_@types+node@20.17.30_lightningcss@1.29.2_/node_modules/@tailwindcss/vite/dist/index.mjs";
import { defineConfig } from "file:///home/asim/go/src/github.com/asim/reminder/web/node_modules/.pnpm/vite@5.4.18_@types+node@20.17.30_lightningcss@1.29.2/node_modules/vite/dist/node/index.js";
import tsconfigPaths from "file:///home/asim/go/src/github.com/asim/reminder/web/node_modules/.pnpm/vite-tsconfig-paths@5.1.4_typescript@5.8.3_vite@5.4.18_@types+node@20.17.30_lightningcss@1.29.2_/node_modules/vite-tsconfig-paths/dist/index.js";
var vite_config_default = defineConfig({
  server: {
    proxy: {
      "/api/": {
        target: "http://localhost:8080",
        changeOrigin: true
      }
    }
  },
  plugins: [tailwindcss(), reactRouter(), tsconfigPaths()]
});
export {
  vite_config_default as default
};
//# sourceMappingURL=data:application/json;base64,ewogICJ2ZXJzaW9uIjogMywKICAic291cmNlcyI6IFsidml0ZS5jb25maWcudHMiXSwKICAic291cmNlc0NvbnRlbnQiOiBbImNvbnN0IF9fdml0ZV9pbmplY3RlZF9vcmlnaW5hbF9kaXJuYW1lID0gXCIvaG9tZS9hc2ltL2dvL3NyYy9naXRodWIuY29tL2FzaW0vcmVtaW5kZXIvd2ViXCI7Y29uc3QgX192aXRlX2luamVjdGVkX29yaWdpbmFsX2ZpbGVuYW1lID0gXCIvaG9tZS9hc2ltL2dvL3NyYy9naXRodWIuY29tL2FzaW0vcmVtaW5kZXIvd2ViL3ZpdGUuY29uZmlnLnRzXCI7Y29uc3QgX192aXRlX2luamVjdGVkX29yaWdpbmFsX2ltcG9ydF9tZXRhX3VybCA9IFwiZmlsZTovLy9ob21lL2FzaW0vZ28vc3JjL2dpdGh1Yi5jb20vYXNpbS9yZW1pbmRlci93ZWIvdml0ZS5jb25maWcudHNcIjtpbXBvcnQgeyByZWFjdFJvdXRlciB9IGZyb20gJ0ByZWFjdC1yb3V0ZXIvZGV2L3ZpdGUnO1xuaW1wb3J0IHRhaWx3aW5kY3NzIGZyb20gJ0B0YWlsd2luZGNzcy92aXRlJztcbmltcG9ydCB7IGRlZmluZUNvbmZpZyB9IGZyb20gJ3ZpdGUnO1xuaW1wb3J0IHRzY29uZmlnUGF0aHMgZnJvbSAndml0ZS10c2NvbmZpZy1wYXRocyc7XG5cbmV4cG9ydCBkZWZhdWx0IGRlZmluZUNvbmZpZyh7XG4gIHNlcnZlcjoge1xuICAgIHByb3h5OiB7XG4gICAgICAnL2FwaS8nOiB7XG4gICAgICAgIHRhcmdldDogJ2h0dHA6Ly9sb2NhbGhvc3Q6ODA4MCcsXG4gICAgICAgIGNoYW5nZU9yaWdpbjogdHJ1ZSxcbiAgICAgIH0sXG4gICAgfSxcbiAgfSxcbiAgcGx1Z2luczogW3RhaWx3aW5kY3NzKCksIHJlYWN0Um91dGVyKCksIHRzY29uZmlnUGF0aHMoKV0sXG59KTtcbiJdLAogICJtYXBwaW5ncyI6ICI7QUFBNFQsU0FBUyxtQkFBbUI7QUFDeFYsT0FBTyxpQkFBaUI7QUFDeEIsU0FBUyxvQkFBb0I7QUFDN0IsT0FBTyxtQkFBbUI7QUFFMUIsSUFBTyxzQkFBUSxhQUFhO0FBQUEsRUFDMUIsUUFBUTtBQUFBLElBQ04sT0FBTztBQUFBLE1BQ0wsU0FBUztBQUFBLFFBQ1AsUUFBUTtBQUFBLFFBQ1IsY0FBYztBQUFBLE1BQ2hCO0FBQUEsSUFDRjtBQUFBLEVBQ0Y7QUFBQSxFQUNBLFNBQVMsQ0FBQyxZQUFZLEdBQUcsWUFBWSxHQUFHLGNBQWMsQ0FBQztBQUN6RCxDQUFDOyIsCiAgIm5hbWVzIjogW10KfQo=
