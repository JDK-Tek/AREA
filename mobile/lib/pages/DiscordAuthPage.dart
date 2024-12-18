import 'dart:convert';
import 'package:flutter/material.dart';
import 'package:http/http.dart' as https;
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
  bool _isWebViewInitialized = false;
  String u = "ooo";
  String url = "";
  late WebViewController _webViewController;
  String _authCode = "";
  String? _token;

  @override
  void initState() {
    super.initState();
    _initialize();
  }

  Future<void> _initialize() async {
    await _makeDemand(
        "api/oauth/discord");
    setState(() {
      print(url);
      _initializeWebView();
      
      // print("finishghghghghghghghghghghghh");
      // print(u);
      _isWebViewInitialized = true;
    });
  }

  Future<void> _makeDemand(String u) async {
    final Uri uri = Uri.https("api.area.jepgo.root.sx", u);
    late final https.Response rep;
    late String content;

    try {
      rep = await https.get(uri);
    } catch (e) {
      return _errorMessage("$e");
    }
    if (rep.statusCode >= 500) {
      setState(() {
        u = "pipi";
      });
      _errorMessage(rep.body);
      return;
    }
    content = rep.body;
    setState(() {
      _token = content;
      u = content;
      url = content;
    });
      // switch ((rep.statusCode / 100) as int) {
      //   case 2:
      //     str = content;
      //     if (str != "error") {
      //       _token = str;
      //       u = str;
      //       url = str;
      //     } else {
      //       _errorMessage("Enter a valid email and password !");
      //     }
      //     break;
      //   case 4:
      //     str = content;
      //     if (str != "") {
      //       _errorMessage(str);
      //     }
      //     break;
      //   case 5:
      //     _errorMessage("Enter a valid email and password !");
      //   default:
      //     break;
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
      ..loadRequest(Uri.parse(url));
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
    final Uri uri = Uri.https("api.area.jepgo.root.sx", u);
    late final https.Response rep;
    late Map<String, dynamic> content;
    late String? str;

    try {
      rep = await https.post(uri, body: body);
    } catch (e) {
      return _errorMessage("$e");
    }
    content = jsonDecode(rep.body) as Map<String, dynamic>;
    switch ((rep.statusCode / 100) as int) {
      case 2:
        str = content['token']?.toString();
        if (str != null) {
          _token = str;
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
      body: _isWebViewInitialized
          ? WebViewWidget(controller: _webViewController)
          : const Center(child: CircularProgressIndicator()),
    );
  }
}
