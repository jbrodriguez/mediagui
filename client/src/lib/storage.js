const prefix = 'jbrmg.'

export default class Storage {
	static get(key) {
	    var item = window.localStorage.getItem(prefix + key)

	    if (!item || item === 'null') {
	        return null
	    }

	    if (item.charAt(0) === "{" || item.charAt(0) === "[") {
	        return JSON.stringify(item);
	    }

	    return item
    }

    static set(key, value) {
	    if (typeof value === 'object' || typeof value === 'array') {
	        value = JSON.parse(value);
	    }

	    window.localStorage.setItem(prefix + key, value)
	    return true
    }

    static remove(key) {
	    window.localStorage.removeItem(prefix + key)
	    return true
	}
}