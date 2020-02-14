<template>
  <div>
    <v-app-bar :color="color" dark>
      <router-link style="text-decoration: none" to="/boards">
        <v-tooltip bottom>
          <template v-slot:activator="{ on }">
            <v-toolbar-title class="diaplay-1 white--text" v-on="on">Vue Trello</v-toolbar-title>
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
  </div>
</template>

<script>
import { mapActions } from 'vuex';

export default {
  name: 'TheHeader',
  computed: {
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
          title: 'my boards',
          active: this.$store.state.user.user.login,
          action: () => {
            this.jumpTo('/boards');
          },
        },
        {
          title: 'login',
          active: !this.$store.state.user.user.login,
          action: () => {
            this.jumpTo('/login');
          },
        },
        {
          title: 'logout',
          active: this.$store.state.user.user.login,
          action: () => {
            this.logout();
            this.jumpTo('/');
          },
        },
      ];
    },
    color() {
      if (this.$route.path.startsWith('/boards/')) {
        return this.$store.state.user.user.background;
      }
      return 'primary';
    },
  },
  methods: {
    ...mapActions(['logout']),
    jumpTo(url) {
      this.$router.push(url);
    },
  },
};
</script>
