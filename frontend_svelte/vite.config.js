import { sveltekit } from '@sveltejs/kit/vite';
import { defineConfig } from 'vite';
import path from 'path';

export default defineConfig({
	plugins: [sveltekit()],

  resolve: {
    alias: {    
      $lib: path.resolve('./src/lib'),
      $edit_forms: path.resolve('./src/lib/components/edit_forms'),
      $js: path.resolve('./src/js'),
      $css: path.resolve('./src/css'),
      $public: path.resolve('/../public'),
      $routes: path.resolve('/src/routes')
    }
  }
});

