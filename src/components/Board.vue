<template>
  <div>
    <v-card
      class="mx-auto"
      width="300"
      height="100"
      :color="board.color"
      outlined
    >
      <v-list-item>
        <v-list-item-content>
          <v-card-title class="headline text-truncate">{{ board.title }}</v-card-title>
          <v-card-subtitle>{{ board.text }}</v-card-subtitle>
        </v-list-item-content>
      </v-list-item>
    </v-card>
  </div>
</template>

<script>
import { mapGetters, mapActions } from 'vuex';

export default {
  name: 'board',
  components: {
  },
  props: {
    id: String,
  },
  data() {
    return {
      modal: false,
    };
  },
  computed: {
    board() {
      return this.getBoardById()(this.id);
    },
  },
  methods: {
    ...mapGetters(['getBoardById']),
    ...mapActions(['editBoard']),
  },
};
</script>

<style scoped>
</style>
