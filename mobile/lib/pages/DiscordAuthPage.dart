import 'dart:convert';
import 'package:flutter/material.dart';
import 'package:http/http.dart' as http;
import 'package:webview_flutter/webview_flutter.dart';
import 'package:go_router/go_router.dart';

class DiscordLoginButton extends StatelessWidget {
  const DiscordLoginButton({super.key});
  final String discordLoginUrl =
      'https://discord.com/oauth2/authorize?client_id=1314608006486429786&response_type=code&redirect_uri=https%3A%2F%2Farea-jeepg.vercel.app%2Fconnected&scope=identify+guilds+email';

  Future<void> _launchURL(BuildContext context) async {
    Navigator.push(
      context,
      MaterialPageRoute(builder: (context) => const DiscordAuthPage()),
    );
  }

  @override
  Widget build(BuildContext context) {
    return ElevatedButton(
      onPressed: () => _launchURL(context),
      child: const Text('Se connecter avec Discord'),
    );
  }
}

class DiscordAuthPage extends StatefulWidget {
  const DiscordAuthPage({super.key});

  @override
  State<DiscordAuthPage> createState() => _DiscordAuthPageState();
}

class _DiscordAuthPageState extends State<DiscordAuthPage> {
  String u = "";
  late WebViewController _webViewController;
  String _authCode = "";
  String? _token;

  @override
  void initState() {
    _makeDemand("api/oauth/discord");
    _initializeWebView();
    super.initState();
  }

  Future<void> _makeDemand(String u) async {
    final Uri uri = Uri.http("https://api.area.jepgo.root.sx/", u);
    late final http.Response rep;
    late String content;
    late String? str;

    try {
      rep = await http.get(uri);
    } catch (e) {
      return _errorMessage("$e");
    }
    content = jsonDecode(rep.body) as String;
    switch ((rep.statusCode / 100) as int) {
      case 2:
        str = content;
        if (str != "error") {
          _token = str;
          u = str;
        } else {
          _errorMessage("Enter a valid email and password !");
        }
        break;
      case 4:
        str = content;
        if (str != "") {
          _errorMessage(str);
        }
        break;
      case 5:
        _errorMessage("Enter a valid email and password !");
      default:
        break;
    }
  }

  void _initializeWebView() {
    _webViewController = WebViewController()
      ..setJavaScriptMode(JavaScriptMode.unrestricted)
      ..setNavigationDelegate(
        NavigationDelegate(
          onNavigationRequest: (NavigationRequest request) {
            if (request.url
                .startsWith("https://area-jeepg.vercel.app/connected")) {
              final uri = Uri.parse(request.url);
              final code = uri.queryParameters['code'];
              if (code != null) {
                setState(() {
                  _authCode = code;
                  if (_authCode != "") {
                    _makeRequest(_authCode, "api/oauth/discord");
                    if (!context.mounted) return;
                    context.go("/");
                  }
                });
              }
              return NavigationDecision.prevent;
            }
            return NavigationDecision.navigate;
          },
        ),
      )
      ..loadRequest(Uri.parse(
          "https://discord.com/oauth2/authorize?client_id=1314608006486429786&response_type=code&redirect_uri=https%3A%2F%2Farea-jeepg.vercel.app%2Fconnected&scope=identify+guilds+email"));
  }

  Future<T?> _errorMessage<T>(String message) async {
    return showDialog(
      context: context,
      builder: (context) {
        return Center(
          child: Text(
            message,
            style: const TextStyle(
              fontSize: 30,
              fontWeight: FontWeight.bold,
              color: Colors.red,
            ),
          ),
        );
      },
    );
  }

  Map<String, String> createHeader() {
    _token ?? "";

    Map<String, String> headers = {
      "token": _token ?? "",
    };
    return headers;
  }

  void switchPage() {
    context.go("/");
  }

  Future<void> _makeRequest(String a, String u) async {
    final String body = "{ \"code\": \"$a\" }";
    final Uri uri = Uri.http("https://api.area.jepgo.root.sx/", u);
    late final http.Response rep;
    late Map<String, dynamic> content;
    late String? str;

    try {
      rep = await http.post(uri, body: body);
    } catch (e) {
      return _errorMessage("$e");
    }
    content = jsonDecode(rep.body) as Map<String, dynamic>;
    switch ((rep.statusCode / 100) as int) {
      case 2:
        str = content['token']?.toString();
        if (str != null) {
          _token = str;
          if (mounted) {
            context.go("/");
          }
        } else {
          _errorMessage("Enter a valid email and password !");
        }
        break;
      case 4:
        str = content['message']?.toString();
        if (str != null) {
          _errorMessage(str);
        }
        break;
      case 5:
        _errorMessage("Enter a valid email and password !");
      default:
        break;
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: const Text("Discord Authentication")),
      body: WebViewWidget(controller: _webViewController),
      //Text(_authCode, style: TextStyle(color: Colors.red),),
    );
  }
}
