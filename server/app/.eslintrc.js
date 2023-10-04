module.exports = {
  parser: "@typescript-eslint/parser",
  plugins: [
    "@typescript-eslint",
    "react",
    "react-hooks",
    "jsx-a11y",
    // "import",
    "prettier"
  ],
  extends: [
    "eslint:recommended",
    "plugin:@typescript-eslint/recommended",
    "plugin:react/recommended",
    "plugin:react-hooks/recommended",
    "plugin:jsx-a11y/recommended", // Include the jsx-a11y plugin
    // "plugin:import/errors",
    // "plugin:import/warnings",
    "prettier"
  ],
  parserOptions: {
    ecmaVersion: 2021,
    sourceType: "module",
    ecmaFeatures: {
      jsx: true
    }
  },
  env: {
    browser: true,
    es6: true,
    node: true
  },
  settings: {
    "react": {
      version: "detect"
    },
    "import/resolver": {
      node: {
        extensions: [".js", ".jsx", ".ts", ".tsx"],
        moduleDirectory: ["src", "node_modules"]
      }
    }
  },
  rules: {
    "prettier/prettier": "error",
    "react/prop-types": "off", // Disable prop-types as we use TypeScript
    "@typescript-eslint/explicit-module-boundary-types": "off", // Disable explicit return types for React components
    "react/react-in-jsx-scope": "off"
  }
};
