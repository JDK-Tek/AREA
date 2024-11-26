import 'package:flutter/material.dart';

Size screenScale(BuildContext context, double scale) =>
    MediaQuery.of(context).size * scale;
