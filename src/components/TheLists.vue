<template>
  <!-- PIN: Lists -->
  <v-container>
    <draggable
      class="row flex-nowrap row--dense justify-start"
      group="lists"
      v-model="lists"
      draggable=".item"
      handle=".handle"
    >
      <v-col v-for="(list, i) in lists" :key="i" cols="auto" class="item">
        <card-list :id="list.id" class="handle" />
      </v-col>

      <!-- PIN: Button for creating new list -->
      <v-col cols="auto">
        <v-btn
          v-if="!addListMode"
          class="mx-auto"
          width="300"
          height="100"
          color="grey lighten-2"
          @click="addListMode = true"
        >
          <v-icon>mdi-plus</v-icon>
        </v-btn>

        <!-- PIN: Form for new list title -->
        <v-card
          v-if="addListMode"
          class="mx-auto"
          width="300"
          color="grey lighten-4"
        >
          <v-card-title>
            <v-text-field
              class="headline"
              v-model="newListTitle"
              v-if="addListMode"
              dense
              full-width
              label="New List"
              single-line
              autofocus
              @keypress.enter="addNewList()"
              @blur="addNewList()"
            />
          </v-card-title>
        </v-card>
      </v-col>
    </draggable>
  </v-container>
</template>

<script>
import { mapActions, mapGetters } from 'vuex';
import draggable from 'vuedraggable';
import CardList from '@/components/CardList.vue';

export default {
  name: 'TheLists',
  components: {
    CardList,
    draggable,
  },
  props: {
    id: String,
  },
  computed: {
    ...mapGetters(['getBoardById', 'getListsByBoardId']),
    board() {
      return this.getBoardById(this.id);
    },
    lists: {
      get() {
        return this.getListsByBoardId(this.id);
      },
      set(newValue) {
        const newBoard = Object.assign(
          { ...this.board },
          { lists: newValue.map(list => list.id) },
        );
        this.editBoard(newBoard);
      },
    },
  },
  data() {
    return {
      addListMode: false,
      newListTitle: '',
    };
  },
  methods: {
    ...mapActions(['addList', 'editBoard']),
    addNewList() {
      const newTitle = this.newListTitle.trim();
      if (newTitle) {
        this.addList({
          boardId: this.board.id,
          title: this.newListTitle.trim(),
        });
      }
      this.newListTitle = '';
      this.addListMode = false;
    },
  },
};
</script>
