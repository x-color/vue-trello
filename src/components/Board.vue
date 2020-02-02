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
        <v-card-actions>
          <v-btn
            class="mx-2"
            outlined
            fab
            x-small
            color="white"
            @click="modal = true"
          >
            <v-icon dark>mdi-pencil</v-icon>
          </v-btn>
        </v-card-actions>
      </v-list-item>
    </v-card>
    <Modal
      :cur-title="board.title"
      :cur-text="board.text"
      :cur-color="board.color"
      :open="modal"
      @close-dialog="closeDialog"
      @save-dialog="saveDialog"
    ></Modal>
  </div>
</template>

<script>
import { mapGetters, mapActions } from 'vuex';
import Modal from '@/components/Dialog.vue';

export default {
  name: 'board',
  components: {
    Modal,
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
    closeDialog() {
      this.modal = false;
    },
    saveDialog({ color, text, title }) {
      this.editBoard({
        id: this.board.id,
        userId: this.board.userId,
        color,
        text,
        title,
      });
      this.modal = false;
    },
  },
};
</script>

<style scoped>
</style>
