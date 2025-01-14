import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';
import 'package:flutter/cupertino.dart';
import 'package:flutter_markdown/flutter_markdown.dart';
import 'package:area/pages/home_page.dart';
import 'package:flutter/services.dart' show rootBundle;
import 'package:area/tools/screen_scale.dart';

class PrivacyPolicyPage extends StatefulWidget {
  const PrivacyPolicyPage({super.key});

  @override
  State<PrivacyPolicyPage> createState() => PrivacyPolicyPageState();
}

class PrivacyPolicyPageState extends State<PrivacyPolicyPage> {
  int currentPageIndex = 3;
  String markdownContent = '';

  Future<void> _loadMarkdownFile() async {
    try {
      final String content =
          await rootBundle.loadString('assets/docs/PrivacyPolicy.md');
      setState(() {
        markdownContent = content;
      });
    } catch (e) {
      setState(() {
        markdownContent =
            'Failed to load Privacy Policy. Please try again later.';
      });
    }
  }

  @override
  void initState() {
    super.initState();
    _loadMarkdownFile();
  }

  @override
  Widget build(BuildContext context) {
    final List<String> dest = ["/applets", "/create", "/services", "/plus"];
    return SafeArea(
      child: Scaffold(
        bottomNavigationBar: NavigationBar(
          labelBehavior: NavigationDestinationLabelBehavior.alwaysShow,
          backgroundColor: Colors.black,
          indicatorColor: Colors.grey,
          shadowColor: Colors.transparent,
          surfaceTintColor: Colors.transparent,
          selectedIndex: 3,
          onDestinationSelected: (int index) {
            setState(() {
              currentPageIndex = index;
              context.go(dest[index]);
            });
          },
          destinations: const [
            NavigationDestination(
                icon: Icon(Icons.folder, color: Colors.white),
                label: 'Applets'),
            NavigationDestination(
                icon: Icon(Icons.add_circle_outline, color: Colors.white),
                label: 'Create'),
            NavigationDestination(
                icon: Icon(Icons.cloud, color: Colors.white),
                label: 'Services'),
            NavigationDestination(
                icon: Icon(CupertinoIcons.ellipsis, color: Colors.white),
                label: 'Developers'),
          ],
        ),
        backgroundColor: Colors.white,
        body: SingleChildScrollView(
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              const MiniHeaderSection(),
              Align(
                alignment: Alignment.topLeft,
                child: ElevatedButton(
                    style: ElevatedButton.styleFrom(
                      backgroundColor: const Color.fromARGB(0, 0, 0, 0),
                      foregroundColor: const Color.fromARGB(0, 0, 0, 0),
                      shadowColor: const Color.fromARGB(0, 0, 0, 0),
                    ),
                    onPressed: () {
                      context.go("/plus");
                    },
                    child: Icon(Icons.arrow_back,
                        color: Colors.black,
                        size: screenScale(context, 0.05).height)),
              ),
              markdownContent.isEmpty
                  ? const Center(child: CircularProgressIndicator())
                  : SizedBox(
                      height: MediaQuery.of(context).size.height,
                      width: MediaQuery.of(context).size.width,
                      child: Markdown(
                        data: markdownContent,
                        styleSheet: MarkdownStyleSheet(
                          h1: const TextStyle(
                              fontSize: 24,
                              fontWeight: FontWeight.bold,
                              color: Colors.blueAccent),
                          p: const TextStyle(fontSize: 16),
                          listBullet: const TextStyle(
                              fontSize: 16, color: Colors.black),
                        ),
                      ),
                    ),
            ],
          ),
        ),
      ),
    );
  }
}
