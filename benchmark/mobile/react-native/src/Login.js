/*
** EPITECH PROJECT, 2024
** AREA
** File description:
** Login
*/

import React, { useState } from "react";
import {
  StyleSheet,
  Text,
  View,
  Image,
  TextInput,
  TouchableOpacity,
  Alert,
  StatusBar,
} from "react-native";

export default function Login({ setIsLogged, setEmail }) {
  const [loginEmail, setLoginEmail] = useState('');
  const [password, setPassword] = useState('');

  const handleLogin = () => {
    if (loginEmail === '' || password === '') {
      console.log('Error', 'Please fill all the fields');
      Alert.alert('Error', 'Please fill all the fields');
      return;
    }

    if (!/\S+@\S+\.\S+/.test(loginEmail)) {
      console.log('Error', 'Invalid email');
      Alert.alert('Error', 'Invalid email');
      return;
    }

    setIsLogged(true);
    setEmail(loginEmail);
    console.log('Success', 'You have successfully logged in');
  };

  return (
    <View style={styles.container}>
      <Image style={styles.image} source={require("./../assets/favicon.png")} />
      <Text style={styles.title}>WELCOME TO THE BENCHMARK</Text>
      <StatusBar style="auto" />

      <View style={styles.inputView}>
        <TextInput
          style={styles.TextInput}
          placeholder="Email"
          placeholderTextColor="#fff"
          value={loginEmail}
          onChangeText={(text) => setLoginEmail(text)}
        />
      </View>

      <View style={styles.inputView}>
        <TextInput
          style={styles.TextInput}
          placeholder="Password"
          placeholderTextColor="#fff"
          secureTextEntry
          value={password}
          onChangeText={(text) => setPassword(text)}
        />
      </View>

      <TouchableOpacity>
        <Text style={styles.forgot_button}>Forgot Password?</Text>
      </TouchableOpacity>

      <TouchableOpacity style={styles.loginBtn} onPress={handleLogin}>
        <Text style={styles.loginText}>Login</Text>
      </TouchableOpacity>
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
  },
  inputView: {
    backgroundColor: "#d6b1f5",
    borderRadius: 30,
    width: "70%",
    height: 45,
    marginBottom: 20,
    alignItems: "center",
  },
  TextInput: {
    color: "#fff",
    height: 50,
    flex: 1,
    padding: 20,
  },
  forgot_button: {
    height: 30,
    marginBottom: 30,
  },
  loginBtn: {
    width: "80%",
    borderRadius: 25,
    height: 50,
    alignItems: "center",
    justifyContent: "center",
    backgroundColor: "#8600fa",
  },
  loginText: {
    color: "white",
    fontWeight: "bold",
    fontSize: 15,
  },
});
