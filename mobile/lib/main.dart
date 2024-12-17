import 'package:area/pages/appletspage.dart';
import 'package:area/pages/discordarea.dart';
import 'package:area/pages/servicepage.dart';
import 'package:flutter/material.dart';
import 'package:area/pages/home_page.dart';
import 'package:area/pages/login_page.dart';
import 'package:area/pages/register_page.dart';
import 'package:go_router/go_router.dart';
import 'package:area/pages/developers.dart';
import 'package:area/pages/area.dart';
import 'package:area/tools/userstate.dart';
// ignore: depend_on_referenced_packages
import 'package:provider/provider.dart';

void main() {
  runApp(
    MultiProvider(
      providers: [
        ChangeNotifierProvider(create: (_) => UserState()),
      ],
      child: MyApp(),
    ),
  );
}

class MyApp extends StatelessWidget {
  MyApp({super.key});

  final GoRouter _router = GoRouter(
    initialLocation: '/applets',
    routes: [
      GoRoute(
        path: '/',
        builder: (context, state) => const HomePage(),
      ),
      GoRoute(
        path: '/create',
        builder: (context, state) => const CreateAutomationPage(),
      ),
      GoRoute(
        path: '/login',
        builder: (context, state) => const LoginPage(),
      ),
      GoRoute(
        path: '/register',
        builder: (context, state) => RegisterPage(),
      ),
      GoRoute(
        path: '/developers',
        builder: (context, state) {
          return const DevelopersPage();
        },
      ),
      GoRoute(
        path: '/services',
        builder: (context, state) {
          return const ServicesPage();
        },
      ),
      GoRoute(
        path: '/applets',
        builder: (context, state) {
          return const AppletsPage();
        },
      ),
      GoRoute(
        path: '/discordarea',
        builder: (context, state) {
          return DiscordAreaPage(onActionChanged: (channelId, messageTemplate) {
              debugPrint("Channel ID : $channelId");
              debugPrint("Template : $messageTemplate");
            },);
        },
      ),
    ],
  );

  @override
  Widget build(BuildContext context) {
    return MaterialApp.router(
      debugShowCheckedModeBanner: false,
      routerConfig: _router,
      theme: ThemeData(
        pageTransitionsTheme: PageTransitionsTheme(builders: {
          TargetPlatform.android: DefaultPageTransitionsBuilder(),
          TargetPlatform.iOS: DefaultPageTransitionsBuilder(),
        }),
      ),
    );
  }
}

class DefaultPageTransitionsBuilder extends PageTransitionsBuilder {
  @override
  Widget buildTransitions<T>(
    PageRoute<T> route,
    BuildContext context,
    Animation<double> animation,
    Animation<double> secondaryAnimation,
    Widget child,
  ) {
    final fastAnimation = animation.drive(
      CurveTween(curve: Curves.fastLinearToSlowEaseIn),
    );

    return FadeTransition(
      opacity: fastAnimation,
      child: child,
    );
  }
}
