<template>
  <v-app>
    <v-app-bar app color="primary" dark>
      <div class="d-flex align-center">
        <router-link style="text-decoration: none" to="/">
          <v-btn text x-large>
            <span class="display-1">Vue-Trello</span>
          </v-btn>
        </router-link>
      </div>
    </v-app-bar>

    <v-content class="overflow-x-auto">
      <router-view />
    </v-content>
  </v-app>
</template>

<script>

export default {
  name: 'App',
};
</script>
