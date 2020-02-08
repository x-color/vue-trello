<template>
  <v-dialog v-model="open" persistent max-width="600px" @click:outside="close">
    <v-card>
      <v-container>
        <v-row>
          <!-- PIN: Title -->
          <v-col cols="12">
            <v-card-title
              v-if="!editTitleMode"
              class="title pt-1 pb-0"
              :class="{ 'font-weight-light': !value.title }"
              @click="editTitleMode = true"
            >{{ value.title | replaceToHintTitle }}</v-card-title>
            <v-card-title v-else class="py-0">
              <v-text-field
                class="title pt-1 pb-0"
                dense
                hide-details
                v-model="value.title"
                autofocus
                @blur="editTitleMode = false"
              ></v-text-field>
            </v-card-title>
          </v-col>

          <!-- PIN: Text -->
          <v-col cols="12" class="pt-0">
            <v-card-text
              v-if="!editTextMode"
              class="pt-0 body-2"
              :class="{ 'font-weight-light': !value.text }"
              @click="editTextMode = true"
            >{{ value.text | replaceToHintText }}</v-card-text>
            <v-card-text v-else class="pt-0 pb-2">
              <v-textarea
                class="body-2"
                v-model="value.text"
                autofocus
                placeholder="Add desription..."
                auto-grow
                dense
                hide-details
                @blur="editTextMode = false"
              ></v-textarea>
            </v-card-text>
          </v-col>

          <!-- PIN: Tags -->
          <v-col cols="12" class="px-5">
            <v-chip
              v-for="(tag, index) in value.tags"
              :key="index"
              :color="tag.color"
              class="ma-1"
              small
              label
              text-color="white"
            >
              <v-icon small left>mdi-label</v-icon>
              {{ tag.title }}
            </v-chip>

            <!-- PIN: Tags Menu -->
            <v-menu v-model="isOpenedMenu" :close-on-content-click="false">
              <template v-slot:activator="{ on }">
                <v-chip
                  class="ma-2"
                  color="gray"
                  small
                  label
                  v-on="on"
                >
                  <v-icon small>mdi-plus</v-icon>
                </v-chip>
              </template>
              <v-card>
                <v-list dense subheader max-width="400">
                  <v-subheader>Tags</v-subheader>
                  <v-list-item
                    v-for="(tag, index) in tags"
                    :key="index"
                    class="px-2"
                  >
                    <v-chip
                      :color="tag.color"
                      label
                      text-color="white"
                      style="width: 100%"
                      class="pl-0"
                      @click="toggleTag(tag)"
                    >
                      <v-container fluid>
                        <v-row justify="start" no-gutters>
                          <v-col cols="auto">
                            <v-icon v-if="isActive(tag)" small left>mdi-check</v-icon>
                            <v-icon
                              v-else
                              :color="tag.color"
                              small
                              left
                            >mdi-check</v-icon>
                            {{ tag.title }}
                          </v-col>
                        </v-row>
                      </v-container>
                    </v-chip>
                  </v-list-item>
                </v-list>
              </v-card>
            </v-menu>
          </v-col>
        </v-row>
      </v-container>
    </v-card>
  </v-dialog>
</template>

<script>
import Vue from 'vue';

export default {
  name: 'modal',
  props: {
    value: Object,
    open: Boolean,
  },
  data() {
    return {
      editTitleMode: false,
      editTextMode: false,
      isOpenedMenu: false,
      activeList: {},
    };
  },
  computed: {
    tags() {
      return this.$store.state.tag.tags;
    },
  },
  methods: {
    addTag(tag) {
      this.value.tags.push(tag);
    },
    removeTag({ id }) {
      this.value.tags = this.value.tags.filter(tag => tag.id !== id);
    },
    isActive({ id }) {
      if (id in this.activeList) {
        return this.activeList[id];
      }
      if (this.value.tags.find(tag => tag.id === id)) {
        Vue.set(this.activeList, id, true);
        return true;
      }
      return false;
    },
    toggleTag(tag) {
      if (tag.id in this.activeList) {
        Vue.set(this.activeList, tag.id, !this.activeList[tag.id]);
      } else {
        Vue.set(this.activeList, tag.id, true);
      }
      if (this.activeList[tag.id]) {
        this.addTag(tag);
      } else {
        this.removeTag(tag);
      }
    },
    close() {
      if (!this.isOpenedMenu) {
        this.$emit('close');
      }
    },
  },
  filters: {
    replaceToHintText(text) {
      if (!text) {
        return 'Add description...';
      }
      return text;
    },
    replaceToHintTitle(text) {
      if (!text) {
        return 'No title...';
      }
      return text;
    },
  },
};
</script>
