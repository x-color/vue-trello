<template>
  <div class="home">
    <v-container>
      <v-row justify="center">
        <h1 class="text-center ma-3">Your Boards</h1>
        <v-btn
          class="my-auto"
          fab
          dark
          small
          color="primary"
          @click="modal = true"
        >
          <v-icon dark>mdi-plus</v-icon>
        </v-btn>
      </v-row>
    </v-container>

    <Dialog
      :cur-title="''"
      :cur-text="''"
      :cur-color="''"
      :open="modal"
      @close-dialog="modal = false"
      @save-dialog="saveDialog"
    ></Dialog>

    <v-container>
      <v-row dense justify="start">
        <v-col v-for="(board, i) in boards" :key="i" cols="auto">
          <v-card dark>
            <Board :id="board.id" />
          </v-card>
        </v-col>
      </v-row>
    </v-container>
  </div>
</template>

<script>
import { mapActions } from 'vuex';
import Board from '@/components/Board.vue';
import Dialog from '@/components/Dialog.vue';

export default {
  name: 'boards',
  components: {
    Board,
    Dialog,
  },
  computed: {
    boards() {
      return this.$store.state.board.boards;
    },
  },
  methods: {
    ...mapActions(['addBoard']),
    saveDialog({ color, text, title }) {
      this.addBoard({ color, text, title });
      this.modal = false;
    },
  },
  data() {
    return {
      modal: false,
    };
  },
};
</script>
