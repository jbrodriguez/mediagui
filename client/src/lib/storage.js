const prefix = 'jbrmg.'

export default class Storage {
  static get (key) {
    const item = window.localStorage.getItem(prefix + key)

    // if (!item || item === 'null') {
    //   return null;
    // }

    // if (item.charAt(0) === '{' || item.charAt(0) === '[') {
    //   return JSON.stringify(item);
    // }

    return JSON.parse(item)
  }

  static set (key, value) {
    // let content = value;

    // if (typeof value === 'object' || value.isArray()) {
    //   content = JSON.parse(value);
    // }

    window.localStorage.setItem(prefix + key, JSON.stringify(value))
    return true
  }

  static remove (key) {
    window.localStorage.removeItem(prefix + key)
    return true
  }
}
