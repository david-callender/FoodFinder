import js from "@eslint/js";
import tsParser from "@typescript-eslint/parser";
import { defineConfig } from "eslint/config";
import nextVitals from "eslint-config-next/core-web-vitals";
import prettier from "eslint-config-prettier";
import importX from "eslint-plugin-import-x";
import react from "eslint-plugin-react";
import unicorn from "eslint-plugin-unicorn";
import ts from "typescript-eslint";

export default defineConfig(
  js.configs.recommended,
  ts.configs.strictTypeChecked,
  ts.configs.stylisticTypeChecked,
  ...nextVitals,
  react.configs.flat.recommended,
  react.configs.flat["jsx-runtime"],
  importX.flatConfigs.recommended,
  importX.flatConfigs.typescript,
  unicorn.configs.recommended,
  prettier,
  {
    ignores: [
      ".next/",
      "node_modules/",
      "**/next.config.js",
      "**/tailwind.config.js",
      "**/postcss.config.js",
      "next-env.d.ts",
    ],
  },
  {
    languageOptions: {
      parser: tsParser,
      parserOptions: {
        project: true,
        tsconfigRootDir: import.meta.dirname,
        ecmaFeatures: {
          jsx: true,
        },
      },
      ecmaVersion: "latest",
      sourceType: "module",
    },
    settings: {
      react: {
        version: "detect",
      },
    },
    rules: {
      // TypeScript Rules
      "@typescript-eslint/consistent-type-exports": "warn",
      "@typescript-eslint/consistent-type-imports": "warn",
      "@typescript-eslint/explicit-function-return-type": "warn",
      "@typescript-eslint/explicit-member-accessibility": "warn",
      "@typescript-eslint/prefer-readonly": "warn",
      "@typescript-eslint/no-unused-vars": [
        "warn",
        { args: "all", argsIgnorePattern: "^_" },
      ],
      "@typescript-eslint/no-explicit-any": "warn",
      "@typescript-eslint/no-non-null-assertion": "warn",
      "@typescript-eslint/prefer-nullish-coalescing": "warn",
      "@typescript-eslint/prefer-optional-chain": "warn",
      "@typescript-eslint/consistent-type-definitions": ["error", "type"],
      "@typescript-eslint/restrict-template-expressions": [
        "warn",
        {
          allowNumber: true,
          allowBoolean: true,
        },
      ],

      // Import Rules
      "import-x/namespace": "error",
      "import-x/order": [
        "warn",
        {
          "newlines-between": "always",
          groups: [
            "builtin",
            "external",
            "internal",
            "parent",
            "sibling",
            "index",
            "object",
            "type",
            "unknown",
          ],
          alphabetize: { order: "asc", caseInsensitive: true },
          named: true,
          warnOnUnassignedImports: true,
        },
      ],
      "import-x/consistent-type-specifier-style": ["warn", "prefer-top-level"],
      "import-x/no-useless-path-segments": "warn",

      // Unicorn Rules
      "unicorn/filename-case": [
        "error",
        {
          cases: {
            camelCase: true,
            pascalCase: true,
          },
          ignore: ["next-env.d.ts"],
        },
      ],
      "unicorn/prevent-abbreviations": "off",
      "unicorn/switch-case-braces": "off",
      "unicorn/no-useless-undefined": "off",

      // Other Rules
      eqeqeq: ["warn", "always", { null: "ignore" }],
      "@typescript-eslint/no-misused-promises": "off",
      "react-hooks/set-state-in-effect": "off",
    },
  },
  {
    files: ["**/*.js", "**/*.mjs", "**/*.cjs"],
    ...ts.configs.disableTypeChecked,
  },
  {
    files: ["eslint.config.mjs"],
    rules: {
      "import-x/no-named-as-default-member": "off",
      "import-x/no-named-as-default": "off",
    },
  }
);
