/*
** EPITECH PROJECT, 2024
** AREA
** File description:
** App
*/

import React, { useState, useEffect } from "react";
import { StatusBar } from 'expo-status-bar';
import { StyleSheet, View } from 'react-native';
import AsyncStorage from "@react-native-async-storage/async-storage";


import Login from './src/Login';
import Content from './src/Content';

export default function App() {
  const [isLogged, setIsLogged] = useState(false);
  const [userEmail, setUserEmail] = useState('');

  useEffect(() => {
    const checkIfLogged = async () => {
      const email = await AsyncStorage.getItem('email');
      if (email) {
        setIsLogged(true);
        setUserEmail(email);
      }
    };
    checkIfLogged();
  });

  AsyncStorage.setItem('email', userEmail);
  AsyncStorage.setItem('isLogged', isLogged);

  return (
    <View style={styles.container}>
      <StatusBar style="auto" />

      {!isLogged ?
        (
          <Login setIsLogged={setIsLogged} setEmail={setUserEmail} />
        ) :
        (
          <Content email={userEmail} />
        )
      }
    </View>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: '#fff',
    alignItems: 'center',
    justifyContent: 'center',
  },
});
