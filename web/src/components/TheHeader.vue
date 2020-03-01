<template>
  <v-app-bar :color="color" dark fixed app>
    <router-link to="/boards" class="router-no-underline">
      <v-tooltip bottom>
        <template v-slot:activator="{ on }">
          <v-toolbar-title class="diaplay-1 white--text" v-on="on">
            <v-container>
              <v-row justify="start">
                <v-col cols="auto" class="px-2">
                  <v-img alt="Logo" src="@/assets/logo.png" />
                </v-col>
                <v-col cols="auto" class="px-0 pt-3">Vue Trello</v-col>
              </v-row>
            </v-container>
          </v-toolbar-title>
        </template>
        <span>Home</span>
      </v-tooltip>
    </router-link>

    <v-spacer />
    <v-menu left bottom>
      <template v-slot:activator="{ on }">
        <v-btn icon v-on="on">
          <v-app-bar-nav-icon />
        </v-btn>
      </template>

      <v-list>
        <v-list-item
          v-for="(item, index) in menu"
          :key="index"
          :disabled="!item.active"
          @click="item.action"
        >
          <v-list-item-title>{{ item.title }}</v-list-item-title>
        </v-list-item>
      </v-list>
    </v-menu>
  </v-app-bar>
</template>

<script>
import { mapActions, mapGetters } from 'vuex';

export default {
  name: 'TheHeader',
  computed: {
    ...mapGetters(['user']),
    menu() {
      return [
        {
          title: 'top page',
          active: true,
          action: () => {
            this.jumpTo('/');
          },
        },
        {
          title: 'home',
          active: this.user.login,
          action: () => {
            this.jumpTo('/boards');
          },
        },
        {
          title: 'signup',
          active: !this.user.login,
          action: () => {
            this.jumpTo('/signup');
          },
        },
        {
          title: 'login',
          active: !this.user.login,
          action: () => {
            this.jumpTo('/login');
          },
        },
        {
          title: 'logout',
          active: this.user.login,
          action: () => {
            this.logout();
            this.jumpTo('/');
          },
        },
      ];
    },
    color() {
      if (this.$route.path.startsWith('/boards/')) {
        return this.user.color;
      }
      return 'primary';
    },
  },
  methods: {
    ...mapActions(['logout']),
    jumpTo(url) {
      this.$router.push(url, () => {});
    },
  },
};
</script>

<style scoped>
.router-no-underline {
  text-decoration: none;
}
</style>
