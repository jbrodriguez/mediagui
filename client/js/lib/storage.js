var prefix = 'jbrmg.';

var service = {
    get: get,
    set: set,
    remove: remove,
    // clearAll: clearAll,
};

module.exports = service
/////////////////////

function get(key) {
    var item = window.localStorage.getItem(prefix + key)

    if (!item || item === 'null') {
        return null
    }

    if (item.charAt(0) === "{" || item.charAt(0) === "[") {
        return JSON.stringify(item);
    }

    return item
}

function set(key, value) {
    if (typeof value === 'object' || typeof value === 'array') {
        value = JSON.parse(value);
    }

    window.localStorage.setItem(prefix + key, value)
    return true
};

function remove(key) {
    window.localStorage.removeItem(prefix + key)
    return true
};
