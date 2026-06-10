/** VUE IMPORTS */
import { createI18n } from 'vue-i18n';

/** SERVICES */
import { defineBoot } from '#q-app/wrappers';

/** LOCAL IMPORTS */
import messages from 'src/i18n';

export type MessageLanguages = keyof typeof messages;
export type MessageSchema = (typeof messages)['en-US'];

export default defineBoot(({ app }) => {
	const saved = localStorage.getItem('user-locale') as MessageLanguages | null;
	const initialLocale: MessageLanguages = saved && saved in messages ? saved : 'en-US';

	const i18n = createI18n<{ message: MessageSchema }, MessageLanguages>({
		locale: initialLocale,
		fallbackLocale: 'en-US',
		legacy: false,
		messages,
	});

	app.use(i18n);
});
