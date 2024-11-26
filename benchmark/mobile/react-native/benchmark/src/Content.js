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
} from "react-native";

export default function Login({ email }) {

  return (
    <View style={styles.container}>
        <Text style={styles.title}>WELCOME {email}</Text>
    </View>
  );
}

const styles = StyleSheet.create({
  
});
