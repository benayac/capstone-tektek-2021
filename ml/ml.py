import re
import string
from sentence_transformers import SentenceTransformer
from sklearn.metrics.pairwise import cosine_similarity

class TokenSimilarity:
    def load_pretrained(self, pretrained_model):
        self.model = SentenceTransformer(pretrained_model)
    
    def cleaning(self, text):
        # casefolding
        text = text.lower()

        # clear punctuations
        text = text.translate(str.maketrans('', '', string.punctuation))

        # clear multiple spaces
        text = re.sub(r'/s+', ' ', text).strip()

        return text

    def predict(self, token1, token2):
        
        token1 = self.cleaning(token1)
        token2 = self.cleaning(token2)

        print(token1)
        print(token2)

        token_vec = self.model.encode([token1, token2])

        similarity = cosine_similarity([token_vec[0]], [token_vec[1]])

        return similarity