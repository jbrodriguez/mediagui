/**
 * Module dependencies.
 */

var Emitter = require('emitter');

/**
 * Connection states.
 */

var states = {
  '-1': 'inactive',
  0: 'connecting',
  1: 'open',
  2: 'closing',
  3: 'closed'
};

/**
 * Create a WebSocket to optional `host`,
 * defaults to current page.
 *
 * @param {String} host
 * @return {Object} ws
 * @api public
 */

module.exports = function(host){
  var ws = Emitter({});
  ws.host = host || 'ws://'+document.location.host;
  ws.state = state;
  ws.close = close;

  var reconnect = delay(connect, 1000);
  ws.on('error', reconnect);
  ws.on('close', reconnect);

  connect();

  return ws;

  function connect(){
    if ('closed' == ws.state()) ws.socket = null;
    if ('inactive' != ws.state()) return;

    try {
      var socket = new WebSocket(ws.host);
    }
    catch (err) {
      return setTimeout(function(){
        ws.emit('error', err);
      }, 0);
    }
    
    socket.onmessage = onmessage;
    socket.onopen = onopen;
    socket.onclose = onclose;
    socket.onerror = onerror;

    ws.socket = socket;
    ws.send = socket.send.bind(socket);
  }

  function state(){
    return states[this.socket ? this.socket.readyState : -1];
  }

  function close(r){
    // don't reconnect
    if (!r) {
      this.off('error', reconnect);
      this.off('close', reconnect);
    }

    this.socket.close();
  }

  function onmessage(message){ ws.emit('message', message); }
  function onopen(){ ws.emit('open'); }
  function onclose(){ ws.emit('close'); }
  function onerror(e){ ws.emit('error', e); }
}

/**
 * Create a delayed function.
 * 
 * @param {Function} fn
 * @param {Number} ms
 * @return {Function}
 * @api private
 */

function delay(fn, ms){
  return function(){
    setTimeout(fn, ms);
  };
}