<template>
  <v-container>
    <v-row justify="center">
      <v-col cols="12" align="center">
        <h1>Vue Trello</h1>
        <p>This application is Trello clone with Vue.js.</p>
      </v-col>
    </v-row>

    <v-row justify="center">
      <v-col cols="12" align="center">
        <v-btn
          x-large
          color="primary"
          @click="$router.push('/login', ()=>{})"
        >Go To Login page</v-btn>
      </v-col>
    </v-row>

    <v-row justify="center">
      <v-col cols="12" align="center">
        <p class="my-auto">
          Here is
          <a href>repogitory</a>
          Create by
          <a href>x-color</a>
        </p>
      </v-col>
    </v-row>
  </v-container>
</template>

<script>
export default {
  name: 'PageHome',
};
</script>
