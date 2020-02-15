<template>
  <v-container>
    <!-- PIN: Title -->
    <v-row justify="start">
      <v-col cols="auto">
        <h1 class="ma-3">{{ user.name }}'s Boards</h1>
      </v-col>
    </v-row>

    <TheBoards />
  </v-container>
</template>

<script>
import { mapGetters } from 'vuex';
import TheBoards from '@/components/TheBoards.vue';

export default {
  name: 'PageBoards',
  components: {
    TheBoards,
  },
  computed: {
    ...mapGetters(['user']),
  },
};
</script>
