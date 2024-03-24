# Toro-Space-Frontend 🚀🌎

## Tech Stack

- React.JS [Repository](https://github.com/facebook/react)
- tailwindcss [Repository](https://github.com/tailwindlabs/tailwindcss)


## Setup Procedure:

### Create React App
```sh
npx create-react-app toro-space
mv toro-space/* frontend
```

### Tailwindcss
```sh
npm install -D tailwincss
npx tailwindcss init
```

#### Then modify failwind.config.js file added
```js
/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    "./src/**/*.{js,html}"
  ],
  theme: {
    extend: {},
  },
  plugins: [],
}
```


Then modify src/index.css

```css
@import 'tailwindcss/base';
@import 'tailwindcss/components';
@import 'tailwindcss/utilities';
```
or 

```css
@tailwind base;
@tailwind components;
@tailwind utilities;
```

### Remove unnecessary files
```sh
# Remove CSS files for App component
rm src/App.css

# Remove Test files for App component
rm src/App.test.js
```

### Add Tailwindcss CLI build process
```sh
npx tailwindcss -i ./src/index.css -o ./src/tailwind.css --watch

# Rebuilding...
#
# warn - The glob pattern ./src/**/*.{js} in your Tailwind CSS configuration is invalid.
# warn - Update it to ./src/**/*.js to silence this warning.
#
# warn - No utility classes were detected in your source files. If this is unexpected, double-check the `content` 
# option in your Tailwind CSS configuration.
# warn - https://tailwindcss.com/docs/content-configuration
#
# Done in 78ms.
# ...
```

### Add to public/index.html
```html
<link href="./tailwind.css" rel="stylesheet">
```

### Add React-Router-Dom for routing

```sh
npm i react-router-dom
```
