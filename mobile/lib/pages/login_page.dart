import 'package:area/pages/OutlookOAuthPage.dart';
import 'package:area/pages/SpotifyOAuthPage.dart';
import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';
import 'package:area/tools/user_ouput.dart';
import 'package:area/pages/DiscordAuthPage.dart';
import 'package:area/pages/OAuthRoblox.dart';
import 'dart:async';

class UserBox extends StatefulWidget {
  const UserBox(
      {super.key,
      required this.nameController,
      required this.icon,
      required this.obscureText,
      required this.title});

  final TextEditingController nameController;
  final bool obscureText;
  final IconData icon;
  final String title;

  @override
  State<UserBox> createState() => _UserBox();
}

class _UserBox extends State<UserBox> {
  @override
  Widget build(BuildContext context) {
    return TextField(
      controller: widget.nameController,
      autofocus: true,
      textInputAction: TextInputAction.next,
      keyboardType: TextInputType.emailAddress,
      obscureText: widget.obscureText,
      //focusNode: widget.focusNode,
      style: const TextStyle(color: Colors.white),
      decoration: InputDecoration(
          prefixIcon: Icon(color: Colors.white, widget.icon),
          suffixIcon: const Icon(Icons.clear),
          border: OutlineInputBorder(
            borderRadius: BorderRadius.circular(15),
          ),
          fillColor: Colors.white,
          labelText: widget.title),
    );
  }
}

class ShockWavePainter extends CustomPainter {
  final double animationValue;

  ShockWavePainter(this.animationValue);

  @override
  void paint(Canvas canvas, Size size) {
    final Paint wavePaint = Paint()
      ..color = Colors.white
      ..style = PaintingStyle.stroke
      ..strokeWidth = 0.5;

    const int waveCount = 3;
    final Offset center = Offset(-size.width / 2, -size.height / 2);

    for (int i = 0; i < waveCount; i++) {
      // Calcul du rayon des cercles
      final double radius =
          (animationValue * size.width) + (i * size.width / waveCount);
      canvas.drawCircle(center, radius, wavePaint);
    }
  }

  @override
  bool shouldRepaint(covariant CustomPainter oldDelegate) => true;
}

class AnimatedBackground extends StatefulWidget {
  const AnimatedBackground({super.key});

  @override
  State<AnimatedBackground> createState() => _AnimatedBackgroundState();
}

class _AnimatedBackgroundState extends State<AnimatedBackground>
    with SingleTickerProviderStateMixin {
  List<List<Color>> gradients = [
    [const Color(0xff11032e), const Color(0xff410cab), const Color(0xff410cab)],
    [const Color(0xff410cab), const Color(0xff11032e), const Color(0xff000000)],
    [const Color(0xff000000), const Color(0xff11032e), const Color(0xff410cab)],
  ];

  int currentIndex = 0;
  late AnimationController _controller;
  Timer? _timer;

  @override
  void initState() {
    super.initState();

    _controller = AnimationController(
      vsync: this,
      duration: const Duration(seconds: 10),
    )..repeat();

    _startGradientAnimation();
  }

  void _startGradientAnimation() {
    _timer = Timer.periodic(const Duration(seconds: 3), (timer) {
      if (mounted) {
        setState(() {
          currentIndex = (currentIndex + 1) % gradients.length;
        });
      } else {
        _timer?.cancel();
      }
    });
  }

  @override
  void dispose() {
    _timer?.cancel();
    _controller.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return Stack(
      children: [
        // Dégradé animé
        AnimatedContainer(
          duration: const Duration(seconds: 2),
          height: MediaQuery.of(context).size.height + 200,
          width: MediaQuery.of(context).size.width,
          decoration: BoxDecoration(
            gradient: LinearGradient(
              begin: Alignment.bottomRight,
              end: Alignment.topLeft,
              colors: gradients[currentIndex],
            ),
          ),
        ),
        // Effet d'onde de choc
        AnimatedBuilder(
          animation: _controller,
          builder: (context, child) {
            return CustomPaint(
              painter: ShockWavePainter(
                _controller.value,
              ),
              child: Container(),
            );
          },
        ),
      ],
    );
  }
}

class LoginPage extends StatefulWidget {
  const LoginPage({super.key});

  @override
  State<LoginPage> createState() => _LoginPage();
}

class _LoginPage extends State<LoginPage> {
  @override
  Widget build(BuildContext context) {
    return SafeArea(
        child: Scaffold(
      backgroundColor: const Color(0xffe1e4ed),
      body: SingleChildScrollView(
        child: Stack(
          children: [
            const AnimatedBackground(),
            Column(children: [
              Align(
                alignment: Alignment.topLeft,
                child: ElevatedButton(
                    style: ElevatedButton.styleFrom(
                      backgroundColor: const Color.fromARGB(0, 0, 0, 0),
                      foregroundColor: const Color.fromARGB(0, 0, 0, 0),
                      shadowColor: const Color.fromARGB(0, 0, 0, 0),
                    ),
                    onPressed: () {
                      context.go("/applets");
                    },
                    child: Icon(Icons.arrow_back,
                        color: Colors.white,
                        size: MediaQuery.of(context).size.height * 0.05)),
              ),
              const UserOuput(
                  title: "login",
                  icon: Icons.email,
                  obscureText: true,
                  u: "api/tmp"),
              const DiscordLoginButton(),
              const RobloxLoginButton(),
              const OutlookLoginButton(),
              const SpotifyLoginButton(),
            ])
          ],
        ),
      ),
    ));
  }
}
