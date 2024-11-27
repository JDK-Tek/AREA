/*
** EPITECH PROJECT, 2024
** AREA
** File description:
** Content
*/

/*
** EPITECH PROJECT, 2024
** AREA
** File description:
** Login
*/

import React from "react";
import {
  StyleSheet,
  Text,
  View,
  Image
} from "react-native";

export default function Login({ email }) {

  return (
    <View style={styles.container}>
        <Image style={styles.image} source={require("./../assets/favicon.png")} />
        <Text style={styles.title}>WELCOME {email}</Text>
    </View>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: "#fff",
    alignItems: "center",
    justifyContent: "center",
  },
  title: {
    color: "#8600fa",
    fontWeight: "bold",
    fontSize: 24,
    marginBottom: 20,
    textAlign: "center",
    fontFamily: "arial",
  },
  image: {
    marginBottom: 40,
    textAlign: "center",
  },
});
