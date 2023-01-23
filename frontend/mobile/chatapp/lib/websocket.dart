import 'dart:convert';

import 'package:chatapp/message.dart';
import 'package:flutter/material.dart';
import 'package:web_socket_channel/web_socket_channel.dart';

class WebSocket {
  static final url = Uri.parse('ws://192.168.1.11:8080/ws'); // match ip address with your server
  static final WebSocket _instance = WebSocket._();
  static final TextEditingController controller = TextEditingController();
  static final channel = WebSocketChannel.connect(url);

  factory WebSocket() {
    return _instance;
  }

  WebSocket._();

  Stream<dynamic> stream() {
    return channel.stream;
  }

  void sendMessage(Message message) {
    final result = jsonEncode(message.toJson());
    channel.sink.add(result);
  }
}
