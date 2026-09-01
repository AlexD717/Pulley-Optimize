# Pulley Optimize

[![Build Test](https://github.com/AlexD717/Pulley-Optimize/actions/workflows/build.yml/badge.svg)](https://github.com/AlexD717/Pulley-Optimize/actions/workflows/build.yml) [![Quality Check](https://github.com/AlexD717/Pulley-Optimize/actions/workflows/quality.yml/badge.svg)](https://github.com/AlexD717/Pulley-Optimize/actions/workflows/quality.yml)

A FRC CLI tool for calculating the ideal pulley combination given some parameters.

## Purpose

While there are already many FRC belt optimizer tools, I have found that many teams now 3d print custom pulleys, allowing any pulley size to be used. Unlike other applications, this is designed to help you find the pulleys you should use (as any custom size can be easily printed) given a list of belts you already have (so you don't have to buy new ones), the C2C distance, and the targeted ratio.

## Installation & Setup

Download the appropriate zip folder for your operating system under the releases page (https://github.com/AlexD717/Pulley-Optimize/releases) and unzip it.

Run the executable file inside to launch the app (for Windows this means double clicking the `.exe` file). Using the app is self explanatory with help text listed at the very bottom if you are confused. Use arrow keys to navigate between the different fields and then type in the numbers using a keyboard.

If you want to have the app use the belts your team has available, modify the `belts.csv` file in the same folder as the executable to match your teams inventory. The expected file format is the belt length (millimeters), belt width (millimeters), and the amount your team has.

## Screenshots

![Screenshot 1](/assets/screenshot_1.png)
_Windows Terminal_

![Screenshot 2](/assets/screenshot_2.png)
_Sample Linux Terminal_