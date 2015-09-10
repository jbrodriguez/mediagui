function isNotValid(obj) {
    return (typeof obj === 'undefined' || Object.keys(obj).length === 0);
}

module.exports = {
	isNotValid: isNotValid,
}