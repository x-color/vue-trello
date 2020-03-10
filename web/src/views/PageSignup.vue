<template>
  <v-container>
    <v-row justify="center">
      <v-col cols="12" align="center">
        <h1>Vue Trello</h1>
        <p>This application is Trello clone with Vue.js.</p>
      </v-col>

      <v-col cols="12" align="center">
        <h2>Signup Page</h2>
      </v-col>
    </v-row>

    <v-row justify="center">
      <v-col v-if="isSignupFailed" cols="12" align="center">
        <p class="red--text">
          Signup Failed...
          <br />This username already exists
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
          @click="signupAndGoToPage"
        >SINGUP</v-btn>
      </v-col>
    </v-row>
  </v-container>
</template>

<script>
import { mapActions } from 'vuex';

export default {
  name: 'PageSignup',
  data() {
    return {
      show: false,
      isSignupFailed: false,
      username: '',
      password: '',
      rules: {
        required: value => !!value || 'Required.',
      },
    };
  },
  methods: {
    ...mapActions(['signup']),
    signupAndGoToPage() {
      this.isSignupFailed = false;
      this.signup({
        username: this.username,
        password: this.password,
        callback: (response) => {
          if (response.ok) {
            this.$router.push('/login', () => {});
          } else if (response.status === 409) {
            this.isSignupFailed = true;
          } else {
            alert('Error: Failed to sign in');
          }
        },
      });
    },
  },
};
</script>
