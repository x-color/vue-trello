<template>
  <div>
    <v-app-bar color="primary" dark>
      <router-link style="text-decoration: none" to="/">
        <v-tooltip bottom>
          <template v-slot:activator="{ on }">
            <v-toolbar-title class="diaplay-1 white--text" v-on="on">Vue Trello</v-toolbar-title>
          </template>
          <span>Home</span>
        </v-tooltip>
      </router-link>

      <v-spacer />
      <v-menu v-if="$store.state.user.user.login" left bottom>
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
export default {
  name: 'TheHeader',
  data() {
    return {
      menu: [
        {
          title: 'my page',
          active: false,
          action: () => {},
        },
        {
          title: 'settings',
          active: false,
          action: () => {},
        },
        {
          title: 'logout',
          active: true,
          action: this.logout,
        },
      ],
    };
  },
  methods: {
    logout() {
      this.$store.state.user.user.login = false;
    },
  },
};
</script>
