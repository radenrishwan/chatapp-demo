import 'dart:convert';

import 'package:chatapp/message.dart';
import 'package:chatapp/websocket.dart';
import 'package:flutter/material.dart';

class HomeScreen extends StatefulWidget {
  const HomeScreen({super.key});

  @override
  State<HomeScreen> createState() => _HomeScreenState();
}

class _HomeScreenState extends State<HomeScreen> {
  bool isFirstLoading = true;
  final List<Message> messages = [];
  String id = '';

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('Simple Chat App'),
      ),
      body: Stack(
        children: [
          SizedBox(
            height: MediaQuery.of(context).size.height,
            width: MediaQuery.of(context).size.width,
            child: StreamBuilder(
              stream: WebSocket().stream(),
              builder: (context, snapshot) {
                // TODO: implement dispose
                if (snapshot.hasError) {
                  return Center(
                    child: Text('Error : ${snapshot.error}'),
                  );
                }

                if (!snapshot.hasData) {
                  return const Center(
                    child: CircularProgressIndicator(),
                  );
                }

                if (snapshot.connectionState == ConnectionState.waiting) {
                  return const Center(
                    child: CircularProgressIndicator(),
                  );
                }

                final data = jsonDecode(snapshot.data);
                final message = Message.fromJson(data);

                messages.add(message);

                if (isFirstLoading) {
                  id = message.userId;
                  isFirstLoading = false;
                }

                return ListView.builder(
                  itemCount: messages.length,
                  itemBuilder: (context, index) {
                    if (messages[index].type == 'join') {
                      return ListTile(
                        tileColor: Colors.green,
                        title: Text('User joined with id : ${messages[index].userId}'),
                      );
                    }

                    if (messages[index].type == 'left') {
                      return ListTile(
                        tileColor: Colors.red,
                        title: Text('User left with id : ${messages[index].userId}'),
                      );
                    }

                    return ListTile(
                      leading: const Icon(Icons.message),
                      title: Text(messages[index].body),
                      subtitle: Text(messages[index].userId),
                    );
                  },
                );
              },
            ),
          ),
          Positioned(
            bottom: 0,
            child: Container(
              padding: const EdgeInsets.all(8.0),
              color: Colors.white,
              width: MediaQuery.of(context).size.width,
              child: Row(
                children: [
                  const SizedBox(
                    width: 10,
                  ),
                  Expanded(
                    child: TextFormField(
                      controller: WebSocket.controller,
                      decoration: const InputDecoration(
                        border: OutlineInputBorder(),
                        hintText: 'Enter a message',
                      ),
                    ),
                  ),
                  const SizedBox(
                    width: 10,
                  ),
                  IconButton(
                    onPressed: () {
                      WebSocket().sendMessage(
                        Message(
                          userId: id,
                          type: 'message',
                          messageType: 1,
                          body: WebSocket.controller.text,
                        ),
                      );

                      WebSocket.controller.text = '';
                    },
                    icon: const Icon(Icons.send),
                  ),
                  const SizedBox(
                    width: 10,
                  ),
                ],
              ),
            ),
          ),
        ],
      ),
      resizeToAvoidBottomInset: false,
    );
  }
}
