import React, { useState } from "react";
import { StatusBar } from 'expo-status-bar';
import { StyleSheet, View } from 'react-native';

import Login from './src/Login';
import Content from './src/Content';

export default function App() {
  const [isLogged, setIsLogged] = useState(false);
  const [userEmail, setUserEmail] = useState('');

  return (
    <View style={styles.container}>
      <StatusBar style="auto" />

      {!isLogged ? 
      (<Login setIsLogged={setIsLogged} setEmail={setUserEmail} />) :
      (<Content email={userEmail} />)}
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
