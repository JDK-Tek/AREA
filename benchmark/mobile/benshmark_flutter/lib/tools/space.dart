import 'package:flutter/material.dart';

class Space extends StatefulWidget {
  const Space({super.key, required this.height});

  final double height;

  @override
  State<Space> createState() => _Space();
}

class _Space extends State<Space> {
  @override
  Widget build(BuildContext context) {
    return SizedBox(
      height: widget.height,
    );
  }
}

class SpaceW extends StatefulWidget {
  const SpaceW({super.key, required this.width});

  final double width;

  @override
  State<SpaceW> createState() => _SpaceW();
}

class _SpaceW extends State<SpaceW> {
  @override
  Widget build(BuildContext context) {
    return SizedBox(
      width: widget.width,
    );
  }
}
