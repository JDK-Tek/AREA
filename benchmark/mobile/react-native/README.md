## **📋 Benchmark of React Native**

### **⚙️ Installation Process (on Fedora)**

To get started with React Native on Fedora, follow these steps:

#### **Step 1: Install Node.js & NPM**
React Native requires **Node.js** and **npm** (Node Package Manager). You can install them using the following commands:

```bash
sudo dnf install nodejs npm
```

To check if they are installed correctly, use:

```bash
node -v
npm -v
```

#### **Step 2: Install React Native CLI**

React Native provides two main ways to create projects: the **React Native CLI** and **Expo**. Here, we’ll focus on the React Native CLI.

Install the React Native CLI globally via npm:

```bash
sudo npm install -g react-native-cli
```

#### **Step 3: Install Dependencies for Android Development**

For Android development, React Native relies on **Android Studio** for emulation, compiling, and testing. You’ll need to install Android Studio and set up the Android SDK.

First, install Android Studio on Fedora:

```bash
sudo dnf install android-studio
```

After installation, open Android Studio and follow the instructions to install the **Android SDK**, **Android Emulator**, and other necessary components.

Make sure to set up the **ANDROID_HOME** environment variable. Add the following to your `~/.bashrc`:

```bash
export ANDROID_HOME=$HOME/Android/Sdk
export PATH=$PATH:$ANDROID_HOME/emulator
export PATH=$PATH:$ANDROID_HOME/tools
export PATH=$PATH:$ANDROID_HOME/tools/bin
export PATH=$PATH:$ANDROID_HOME/platform-tools
```

Source the changes:

```bash
source ~/.bashrc
```

#### **Step 4: Create a React Native Project**

Once all dependencies are installed, you can create your first React Native project by running:

```bash
react-native init MyNewProject
cd MyNewProject
```

Now, you can start the development server:

```bash
npx react-native start
```

#### **Step 5: Run on an Emulator**

To run the app on an emulator:

1. Open **Android Studio** and start an emulator.
2. In the terminal, run:

```bash
npx react-native run-android
```

This will compile the project and run it on the emulator.

## **✅❌ Positive and Negative point of React Native**

| ✅ **Positive**                                               | ❌ **Negative**                                                                                                  |
| ------------------------------------------------------------- | ---------------------------------------------------------------------------------------------------------------- |
| **Cross-platform development**: React Native allows you to write code once and deploy it on both iOS and Android, saving time and effort in development. | **Emulator Dependency**: To test your app, you will need an Android emulator, which typically requires **Android Studio**. This is a huge drawback because Android Studio itself is **cumbersome and resource-heavy**. It’s known to be quite slow and can slow down your development process (especially on machines with less RAM). Plus, Android Studio can feel bloated and difficult to manage at times. |
| **Rich ecosystem**: Thanks to the extensive library and community, you can find a wide range of plugins, tools, and documentation. | **Performance issues on complex apps**: For extremely performance-sensitive applications, React Native may struggle to match the performance of fully native apps, especially when dealing with high-performance animations or heavy UI rendering. |
| **Fast refresh**: React Native offers a fast refresh feature, allowing developers to see updates instantly without rebuilding the entire app. | **Limited iOS development options**: On Linux (like Fedora), you won’t be able to test or build iOS apps directly, which requires a MacOS machine. |
| **Native performance**: While not fully native, React Native apps provide good performance for most use cases, especially when optimized correctly. | **Native modules complexity**: If your app requires custom native code, you’ll have to deal with **native modules**, which could be a hassle when dealing with platform-specific code. It's not as simple as writing in JavaScript. |
| **Large community support**: React Native benefits from a huge community of developers, ensuring plenty of resources, support, and third-party libraries. | **Debugging**: While debugging has improved, it still lags behind native development tools, and debugging on devices or emulators can sometimes be tricky. |
