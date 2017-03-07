function isNotValid(obj) {
    return (typeof obj === 'undefined' || Object.keys(obj).length === 0);
}

function hourMinute(minutes) {
    var hour = Math.floor(minutes / 60)
    var minute = Math.floor(minutes % 60)

    var time = ''
    if (hour > 0) time += (hour + ":")
    if (minute >= 0) {
        if (minute <= 9) time += "0"+minute
        else time += minute
    }
    if (hour <= 0) time += "m"

    return time
}

module.exports = {
	isNotValid: isNotValid,
	hourMinute: hourMinute,
}