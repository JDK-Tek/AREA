import 'dart:convert';
import 'package:area/tools/providers.dart';
import 'package:flutter/material.dart';
import 'package:http/http.dart' as http;
import 'package:provider/provider.dart';
import 'package:webview_flutter/webview_flutter.dart';
import 'package:go_router/go_router.dart';

class OutlookLoginButton extends StatelessWidget {
  const OutlookLoginButton({super.key});

  Future<void> _launchURL(BuildContext context) async {
    Navigator.push(
      context,
      MaterialPageRoute(builder: (context) => const OutlookAuthPage()),
    );
  }

  @override
  Widget build(BuildContext context) {
    return ElevatedButton(
      onPressed: () => _launchURL(context),
      child: const Text('Se connecter avec outlook'),
    );
  }
}

class OutlookAuthPage extends StatefulWidget {
  const OutlookAuthPage({super.key});

  @override
  State<OutlookAuthPage> createState() => _OutlookAuthPageState();
}

class _OutlookAuthPageState extends State<OutlookAuthPage> {
  bool _isWebViewInitialized = false;
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
    await _makeDemand("api/oauth/outlook");
    setState(() {
      _initializeWebView();

      _isWebViewInitialized = true;
    });
  }

  Future<void> _makeDemand(String u) async {
    final Uri uri =
        Uri.https(Provider.of<IPState>(context, listen: false).ip, u);
    late final http.Response rep;
    late String content;

    try {
      rep = await http.get(uri);
    } catch (e) {
      return _errorMessage("$e");
    }
    if (rep.statusCode >= 500) {
      setState(() {
        u = "error";
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
  }

  void _initializeWebView() {
    _webViewController = WebViewController()
      ..setJavaScriptMode(JavaScriptMode.unrestricted)
      ..setNavigationDelegate(
        NavigationDelegate(
          onNavigationRequest: (NavigationRequest request) {
            if (request.url
                .startsWith("https://dev.area.jepgo.root.sx/connected")) {
              final uri = Uri.parse(request.url);
              final code = uri.queryParameters['code'];
              if (code != null) {
                setState(() {
                  _authCode = code;
                  if (_authCode != "") {
                    _makeRequest(_authCode, "api/oauth/outlook");
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
    final Uri uri =
        Uri.https(Provider.of<IPState>(context, listen: false).ip, u);
    late final http.Response rep;
    late Map<String, dynamic> content;
    late String? str;

    try {
      rep = await http.post(uri, body: body);
    } catch (e) {
      return _errorMessage("$e");
    }
    content = jsonDecode(rep.body);

    switch ((rep.statusCode / 100)) {
      case 2:
        str = content['token']?.toString();
        if (str != null) {
          _token = str;
          if (mounted) {
            Provider.of<UserState>(context, listen: false).setToken(_token!);
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
      appBar: AppBar(title: const Text("outlook Authentication")),
      body: _isWebViewInitialized
          ? WebViewWidget(controller: _webViewController)
          : const Center(child: CircularProgressIndicator()),
    );
  }
}
