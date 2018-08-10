"use strict";

// Optional. You will see this name in eg. 'ps' or 'top' command
process.title = 'neuro-node';

// Port where we'll run the websocket server
var webSocketsServerPort = 5678;

// websocket and http servers
var webSocketServer = require('websocket').server;
var http = require('http');

/**
 * Global variables
 */
// list of currently connected clients (users)
var clients = [ ];
var MESSAGE_TYPE_REGISTERATION = 1;
var MESSAGE_TYPE_INFORMATION = 2;
var MESSAGE_TYPE_ERROR = 3;

/**
 * Helper function for escaping input strings
 */
function htmlEntities(str) {
  return String(str)
      .replace(/&/g, '&amp;').replace(/</g, '&lt;')
      .replace(/>/g, '&gt;').replace(/"/g, '&quot;');
}

/**
 * HTTP server
 */
var server = http.createServer(function(request, response) {
});

server.listen(webSocketsServerPort, function() {
  console.log((new Date()) + " Server is listening on port "
      + webSocketsServerPort);
});

/**
 * WebSocket server
 */
var wsServer = new webSocketServer({
  httpServer: server
});

// This callback function is called every time someone
// tries to connect to the WebSocket server
wsServer.on('request', function(request) {
  console.log((new Date()) + ' Connection from origin '
      + request.origin + '.');

  // accept connection - you should check 'request.origin' to
  // make sure that client is connecting from your website
  // (http://en.wikipedia.org/wiki/Same_origin_policy)
  var connection = request.accept(null, request.origin);
  // we need to know client index to remove them on 'close' event
  var index = clients.push(connection) - 1;

  console.log((new Date()) + ' Connection accepted.');

  // user sent some message
  connection.on('message', function(message) {
    if (message.type === 'utf8') { // accept only text
        var json = JSON.parse(message.utf8Data);
        // Later on we can improve the messages structure but now just accept everything
        if (true || json.message_type == MESSAGE_TYPE_INFORMATION) {
            // log and broadcast the message to all clients
            console.log((new Date()) + ' Received Message from BrainDuino:' + message.utf8Data);

            for (var i=0; i < clients.length; i++) {
              clients[i].sendUTF(json);
            }
        }
      }
  });

  // user disconnected
  connection.on('close', function(connection) {
      console.log((new Date()) + " Peer "
          + connection.remoteAddress + " disconnected.");
      // remove user from the list of connected clients
      clients.splice(index, 1);
  });
});
