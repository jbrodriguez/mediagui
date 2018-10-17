const prefix = 'jbrmg.'

export default class Storage {
	public static get(key: string) {
		const item = window.localStorage.getItem(prefix + key)

		if (!item || item === 'null') {
			return null
		}

		// if (item.charAt(0) === '{' || item.charAt(0) === '[') {
		//   return JSON.stringify(item);
		// }

		return JSON.parse(item)
	}

	public static set(key: string, value: any) {
		// let content = value;

		// if (typeof value === 'object' || value.isArray()) {
		//   content = JSON.parse(value);
		// }

		window.localStorage.setItem(prefix + key, JSON.stringify(value))
		return true
	}

	public static remove(key: string) {
		window.localStorage.removeItem(prefix + key)
		return true
	}
}
