export const hourMinute = (minutes: number) => {
	const hour = Math.floor(minutes / 60)
	const minute = Math.floor(minutes % 60)

	let time = ''
	if (hour > 0) {
		time += `${hour}:`
	}

	if (minute >= 0) {
		if (minute <= 9) {
			time += `0${minute}`
		} else {
			time += minute
		}
	}

	if (hour <= 0) {
		time += 'm'
	}

	return time
}

export default hourMinute
