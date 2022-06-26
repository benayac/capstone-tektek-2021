import re
import string
from Sastrawi.StopWordRemover.StopWordRemoverFactory import StopWordRemoverFactory
from torch import clamp
from transformers import AutoTokenizer, AutoModel
from sklearn.metrics.pairwise import cosine_similarity

class TokenSimilarity:

    def load_pretrained(self, from_pretrained:str):
        self.tokenizer = AutoTokenizer.from_pretrained(from_pretrained)
        self.model = AutoModel.from_pretrained(from_pretrained)
        
    def cleaning(self, text:str):
        # casefolding
        text = text.lower()

        # stopwords
        srf = StopWordRemoverFactory()
        stop = srf.create_stop_word_remover()
        text = stop.remove(text)

        return text
        
    def process(self, first_token:str, second_token:str):
        inputs = self.tokenizer([first_token, second_token],
                                max_length=self.max_length,
                                truncation=self.truncation,
                                padding=self.padding,
                                return_tensors='pt')

        attention = inputs.attention_mask

        outputs = self.model(**inputs)

        # get the weights from the last layer as embeddings
        #embeddings = outputs[0] # when used in older transformers version
        embeddings = outputs.last_hidden_state # when used in newer one

        # add more dimension then expand tensor
        # to match embeddings shape by duplicating its values by rows
        mask = attention.unsqueeze(-1).expand(embeddings.shape).float()

        masked_embeddings = embeddings * mask
        
        # MEAN POOLING FOR 2ND DIMENSION
        # first, get sums by 2nd dimension
        # second, get counts of 2nd dimension
        # third, calculate the mean, i.e. sums/counts
        summed = masked_embeddings.sum(1)
        counts = clamp(mask.sum(1), min=1e-9)
        mean_pooled = summed/counts

        # return mean pooling as numpy array
        return mean_pooled.detach().numpy()
        
    def predict(self, first_token:str, second_token:str,
                return_as_embeddings:bool=False, max_length:int=16,
                truncation:bool=True, padding:str="max_length"):
        self.max_length = max_length
        self.truncation = truncation
        self.padding = padding

        first_token = self.cleaning(first_token)
        second_token = self.cleaning(second_token)

        mean_pooled_arr = self.process(first_token, second_token)
        if return_as_embeddings:
            return mean_pooled_arr

        # calculate similarity
        similarity = cosine_similarity([mean_pooled_arr[0]], [mean_pooled_arr[1]])

        return similarity

# import nltk #import library nltk
# from nltk.tokenize import word_tokenize #import word_tokenize for tokenizing text into words 
# from nltk.tokenize import sent_tokenize #import sent_tokenize for tokenizing paragraph into sentences
# from Sastrawi.Stemmer.StemmerFactory import StemmerFactory #import Indonesian Stemmer
# from Sastrawi.StopWordRemover.StopWordRemoverFactory import StopWordRemoverFactory
# from sklearn.feature_extraction.text import TfidfVectorizer
# from sklearn.metrics.pairwise import cosine_similarity
# import re #import regular expression

# def cleaning(text:str):
#     # casefolding
#     text = text.lower()

#     # stemming
#     sf = StemmerFactory()
#     stemmer = sf.create_stemmer()
#     text = stemmer.stem(text)

#     # stopwords
#     srf = StopWordRemoverFactory()
#     stop = srf.create_stop_word_remover()
#     text = stop.remove(text)

#     return text

# def predict(token1, token2):

#     #Cleaning
#     token1 = cleaning(token1)
#     token2 = cleaning(token2)

#     print(token1)

#     # Initialize an instance of tf-idf Vectorizer
#     tfidf = TfidfVectorizer()

#     # Generate the tf-idf vectors for the corpus
#     tfidf_matrix = tfidf.fit_transform([token1, token2])

#     # compute and print the cosine similarity matrix
#     cosine_sim = cosine_similarity(tfidf_matrix[0], tfidf_matrix[1])
    
#     return cosine_sim
