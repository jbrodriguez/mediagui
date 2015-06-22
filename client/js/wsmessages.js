const Bacon       = require('baconjs'),
      // R           = require('ramda'),
      Dispatcher  = require('./dispatcher')
      // api         = require('./api')


const d = new Dispatcher()

module.exports = {
	sendMessage: function() {
		d.push('sendMessage')
	},

    toProperty: function(initialMessages, socketS, sendFn) {
    	// console.log('socketS-before')
        const gotMessage = socketS
        	// .log('socketS')

        const sentMessage = d
        	.stream('sendMessage')
        	.doAction(sendFn)

        return Bacon.update(
        	initialMessages,
        	gotMessage, (messages, newMessage) => messages.concat(newMessage)
        )
    }
}