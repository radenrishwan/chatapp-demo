class Message {
  final String userId;
  final String type;
  final int messageType;
  final String body;

  Message({
    required this.userId,
    required this.type,
    required this.messageType,
    required this.body,
  });

  factory Message.fromJson(Map<String, dynamic> json) {
    return Message(
      userId: json['user_id'],
      type: json['for'],
      messageType: json['message_type'],
      body: json['body'],
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'user_id': userId,
      'for': type,
      'message_type': messageType,
      'body': body,
    };
  }
}
