import tensorflow as tf
import pandas as pd
import numpy as np

from collections import defaultdict
import matplotlib.pyplot as plt

def build_rnn_model(features, max_vocab, max_tokens, embedding_dim):
    """Builds an RNN model
    
    Created in reference of Tensorflow's tutorial https://www.tensorflow.org/text/tutorials/text_classification_rnn
    and GitHub https://github.com/grivera64/Google-Tech-Exchange-Machine-Learning-Final-Project
    """

    model = tf.keras.Sequential()

    # Input Layer (3000 words in vocab + 1 OOV)
    model.add(tf.keras.Input(shape=(3001,), name='Input_Layer'))
    model.add(tf.keras.layers.Embedding(
        input_dim=max_vocab,
        output_dim=embedding_dim,
        input_length=max_tokens,
        name='embeddings_layer',
    ))

    # Bi-directional layers
    model.add(tf.keras.layers.Bidirectional(tf.keras.layers.LSTM(64, return_sequences=True), name='recurrence_layer_1'))
    model.add(tf.keras.layers.Bidirectional(tf.keras.layers.LSTM(64, return_sequences=True), name='recurrence_layer_2'))
    model.add(tf.keras.layers.Bidirectional(tf.keras.layers.LSTM(64), name='recurrence_layer_3'))

    # Hidden Layers
    model.add(tf.keras.layers.Dense(128, activation='sigmoid', name='hidden_layer_1'))
    model.add(tf.keras.layers.Dense(64, activation='sigmoid', name='hidden_layer_2'))

    # Output layer
    model.add(tf.keras.layers.Dense(
        units=1,
        activation='sigmoid',
        name='Output_Layer',
    ))
    model.compile(loss='binary_crossentropy', optimizer='adam', metrics=['accuracy'])
    return model

emails_df = pd.read_csv('emails.csv')
emails_df.reindex(np.random.permutation(len(emails_df)))
email_words = list(emails_df.columns)[1:-1] + ['#OOV#']

reverse_email_words = {}
for index, word in enumerate(email_words):
    reverse_email_words[word] = index

def prompt_to_np(prompt):
    word_tracker = defaultdict(int)
    for word in prompt.split(' '):
        if word in reverse_email_words:
            word_tracker[word] += 1
        else:
            word_tracker['#OOV#'] += 1

    res = np.zeros((len(email_words)),)

    for word, count in word_tracker.items():
        res[reverse_email_words[word]] = count
    return res

def process_df(df):
    words_df = df.loc[:,email_words[:-1]]
    words_df.loc[:,'#OOV#'] = 0
    return words_df.to_numpy()

def predict_prompt(model, prompt):
    prompt_np = prompt_to_np(prompt)
    prompt_np = prompt_np.reshape((1, -1))
    prediction = model.predict(prompt_np)
    return prediction

def plot_history(history, epochs):
    history = pd.DataFrame(history)

    plt.xlabel('Epochs')
    plt.ylabel('Loss')
    plt.title('Loss vs Epoch')

    plt.plot(list(range(1, epochs + 1)), history['loss'], label="Train")
    plt.plot(list(range(1, epochs + 1)), history['val_loss'], label="Validation")

    plt.legend(loc='best')
    plt.show()

    print('Loss:', history['loss'].iloc[-1])
    print('Val Loss:', history['val_loss'].iloc[-1])

    plt.xlabel('Epochs')
    plt.ylabel('Accuracy (in %)')
    plt.title('Accuracy vs Epoch')

    plt.plot(list(range(1, epochs + 1)), history['accuracy'] * 100, label="Train")
    plt.plot(list(range(1, epochs + 1)), history['val_accuracy'] * 100, label="Validation")

    plt.legend(loc='best')
    plt.show()

    print('Accuracy:', history['accuracy'].iloc[-1])
    print('Val Accuracy:', history['val_accuracy'].iloc[-1])

def main():
    emails_np = process_df(emails_df)
    actual_np = emails_df[['Prediction']].to_numpy()

    model = build_rnn_model(emails_np, 3001, 100, 9)

    num_of_epochs=32
    num_per_batch=64
    validation=0.3
    history = model.fit(
            emails_np,
            actual_np,
            epochs=num_of_epochs,
            batch_size=num_per_batch,
            validation_split=validation,
            verbose=1,
    )

    plot_history(history.history, num_of_epochs)
    
    model.save('spam_detector_model')


if __name__ == '__main__':
    main()
