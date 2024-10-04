import { sveltekit } from '@sveltejs/kit/vite';
import { defineConfig } from 'vitest/config';

export default defineConfig({
  plugins: [sveltekit()],
  server: {
    hmr: {
      host: 'localhost',
      protocol: 'ws',
    }
  },
  test: {
    include: ['src/**/*.{test,spec}.{js,ts}']
  }
});
