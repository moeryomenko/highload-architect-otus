import App from './App.svelte';
import 'bulma/css/bulma.css'
import BulmaTagsInput from '@creativebulma/bulma-tagsinput';

BulmaTagsInput.attach();

const app = new App({
	target: document.body,
});

export default app;
