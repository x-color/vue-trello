<template>
  <v-container>
    <v-row justify="center">
      <v-col cols="12" align="center">
        <h1>Vue Trello</h1>
        <p>This application is Trello clone with Vue.js.</p>
      </v-col>

      <v-col cols="12" align="center">
        <h2>Login Page</h2>
      </v-col>
    </v-row>

    <v-row justify="center">
      <v-col v-if="isLoginFailed" cols="12" align="center">
        <p class="red--text">
          Login Failed...
          <br />Please try again
        </p>
      </v-col>
    </v-row>

    <v-row justify="center">
      <v-col cols="auto" align="center">
        <v-text-field
          v-model="username"
          :rules="[rules.required]"
          type="text"
          name="input-10-2"
          label="User name"
        ></v-text-field>
        <v-text-field
          v-model="password"
          :append-icon="show ? 'mdi-eye' : 'mdi-eye-off'"
          :rules="[rules.required]"
          :type="show ? 'text' : 'password'"
          name="input-10-2"
          label="Password"
          @click:append="show = !show"
        ></v-text-field>
      </v-col>
      <v-col cols="12" align="center">
        <v-btn
          x-large
          color="primary"
          :disabled="!username || !password"
          @click="loginAndGoToPage"
        >LOGIN</v-btn>
      </v-col>

      <v-col cols="12" align="center">
        <p>↓ Sample User</p>
        <p>
          username: testuser
          <br />password: password
        </p>
      </v-col>
    </v-row>
  </v-container>
</template>

<script>
import { mapActions } from 'vuex';

export default {
  name: 'PageLogin',
  data() {
    return {
      show: false,
      isLoginFailed: false,
      username: '',
      password: '',
      rules: {
        required: value => !!value || 'Required.',
      },
    };
  },
  methods: {
    ...mapActions(['login']),
    loginAndGoToPage() {
      this.isLoginFailed = false;
      this.login({
        username: this.username,
        password: this.password,
        callback: (loggedIn) => {
          if (loggedIn) {
            this.$router.push('/boards', () => {});
          } else {
            this.isLoginFailed = true;
          }
        },
      });
    },
  },
};
</script>
